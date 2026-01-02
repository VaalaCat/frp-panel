package upgrade

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/VaalaCat/frp-panel/utils"
	"github.com/VaalaCat/frp-panel/utils/logger"
	"github.com/kardianos/service"
)

const upgraderServiceName = "frpp-upgrader"

type StartResult struct {
	Dispatched      bool
	PlanPath        string
	UpgraderService string
}

// Start 执行升级（非 Windows：直接替换不影响当前进程；Windows：启动 worker 完成替换/服务控制）
func Start(ctx context.Context, opt Options) error {
	_, err := StartWithResult(ctx, opt)
	return err
}

func StartWithResult(ctx context.Context, opt Options) (StartResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	// 关键诊断信息：帮助确认“到底执行的是哪个 frp-panel 二进制”
	exePath, _ := os.Executable()
	realExe := exePath
	if rp, err := filepath.EvalSymlinks(exePath); err == nil && len(rp) > 0 {
		realExe = rp
	}
	if abs, err := filepath.Abs(realExe); err == nil {
		realExe = abs
	}

	target, err := resolveTargetPath(opt.TargetPath)
	if err != nil {
		return StartResult{}, err
	}
	opt.TargetPath = target
	if len(strings.TrimSpace(opt.Version)) == 0 {
		opt.Version = "latest"
	}
	if len(opt.WorkDir) == 0 {
		opt.WorkDir = defaultWorkDir()
	}

	unlock, err := lock(opt.WorkDir)
	if err != nil {
		return StartResult{}, err
	}
	defer unlock()

	if err := utils.EnsureDirectoryExists(opt.TargetPath); err != nil {
		return StartResult{}, fmt.Errorf("ensure target directory failed: %w", err)
	}

	downloadURL, err := buildDownloadURL(opt)
	if err != nil {
		return StartResult{}, err
	}

	logger.Logger(ctx).Infof("upgrade: executable=%s, target=%s, restart_service=%v, service_name=%s",
		realExe, opt.TargetPath, opt.RestartService, strings.TrimSpace(opt.ServiceName))
	logger.Logger(ctx).Infof("upgrade: downloading version [%s], url: %s", opt.Version, downloadURL)
	tmpPath, err := utils.DownloadFile(ctx, downloadURL, strings.TrimSpace(opt.HTTPProxy))
	if err != nil {
		return StartResult{}, fmt.Errorf("download failed: %w", err)
	}

	// stage 到目标目录附近，避免跨文件系统 rename 问题
	staged := stagePathForTarget(opt.TargetPath)
	_ = os.Remove(staged)
	if err := copyFile(tmpPath, staged, 0755); err != nil {
		return StartResult{}, fmt.Errorf("stage new binary failed: %w", err)
	}
	_ = os.Chmod(staged, 0755)

	if err := verifyBinary(staged); err != nil {
		return StartResult{}, err
	}

	// Linux 场景：避免“远程升级递归依赖”（stop/restart frpp 会杀掉同 unit/cgroup 下的升级进程）
	// 做法：写入固定 plan.json，然后启动独立的 upgrader service 去 stop→替换→start。
	if runtime.GOOS == "linux" && opt.RestartService && len(strings.TrimSpace(opt.ServiceName)) > 0 {
		planPath, err := writePlan(opt.WorkDir, opt)
		if err != nil {
			return StartResult{}, err
		}
		logger.Logger(ctx).Infof("upgrade: plan created at %s, dispatching upgrader service: %s", planPath, upgraderServiceName)

		if err := ensureUpgraderService(ctx, planPath); err != nil {
			return StartResult{}, err
		}
		// 异步：此处返回后 remoteshell 可以立刻得到响应；真正 stop/restart 将由 upgrader 完成
		return StartResult{Dispatched: true, PlanPath: planPath, UpgraderService: upgraderServiceName}, nil
	}

	// Windows：无法覆盖正在运行的 exe，因此用独立 worker 来做服务 stop->replace->start
	if runtime.GOOS == "windows" {
		planPath, err := writePlan(opt.WorkDir, opt)
		if err != nil {
			return StartResult{}, err
		}
		// 让 worker 复用已 stage 的文件：约定 staged 固定为 targetPath+".new"
		if err := spawnWorker(mustExecutablePath(), planPath); err != nil {
			return StartResult{}, fmt.Errorf("start upgrade worker failed: %w", err)
		}
		logger.Logger(ctx).Info("upgrade worker started (windows). it will stop/replace/start service in background if configured")
		return StartResult{Dispatched: true, PlanPath: planPath, UpgraderService: ""}, nil
	}

	// 非 Windows：当前进程可以继续运行，替换不会影响当前运行实例
	var backupPath string
	if opt.Backup {
		backupPath, err = backupExisting(opt.TargetPath)
		if err != nil {
			return StartResult{}, err
		}
	}

	if err := replaceFile(staged, opt.TargetPath); err != nil {
		// 回滚
		if len(backupPath) > 0 {
			_ = replaceFile(backupPath, opt.TargetPath)
		}
		return StartResult{}, fmt.Errorf("replace executable failed: %w", err)
	}

	logger.Logger(ctx).Infof("upgrade: binary replaced successfully: %s", opt.TargetPath)

	if opt.RestartService && len(strings.TrimSpace(opt.ServiceName)) > 0 {
		logger.Logger(ctx).Infof("upgrade: restarting service: %s", opt.ServiceName)
		// 参考 cmd/frpp/shared/cmd.go：使用 utils.ControlSystemService
		if err := utils.ControlSystemService(opt.ServiceName, opt.ServiceArgs, "restart", func() {}); err != nil {
			// 二进制已经替换成功，这里的重启失败不应导致整体 upgrade 失败（尤其是非 root 场景）
			logger.Logger(ctx).WithError(err).Warnf("restart service failed, please restart manually or run with sudo: %s", opt.ServiceName)
			return StartResult{}, nil
		}
	}

	return StartResult{Dispatched: false, PlanPath: "", UpgraderService: ""}, nil
}

func ensureUpgraderService(ctx context.Context, planPath string) error {
	// 统一使用 kardianos/service：兼容非 systemd 的 Linux（upstart/sysv/openrc）
	// 对 systemd：通过 SystemdScript/Restart 选项，避免 Restart=always 导致无限重启，
	// 并用 ConditionPathExists 防止 enable 后开机自启误触发升级。

	args := []string{"__upgrade-worker", "--plan", planPath}

	opts := service.KeyValue{
		// 默认 systemd 脚本会 Restart=always 且 install 会 enable，这会导致无限重启。
		// 我们定制为 oneshot + Restart=no + ConditionPathExists(plan)。
		"SystemdScript": systemdUpgraderScript(planPath),
		"Restart":       "no",
	}

	// 修复/覆盖旧 unit：stop + uninstall（忽略错误）后 install + start
	_ = utils.ControlSystemServiceWithOptions(upgraderServiceName, args, "stop", func() {}, opts)
	_ = utils.ControlSystemServiceWithOptions(upgraderServiceName, args, "uninstall", func() {}, opts)
	if err := utils.ControlSystemServiceWithOptions(upgraderServiceName, args, "install", func() {}, opts); err != nil {
		// 如果 install 因为已存在等原因失败，再尝试直接 start
		logger.Logger(ctx).WithError(err).Warn("upgrade: upgrader install failed, try start directly")
	}
	return utils.ControlSystemServiceWithOptions(upgraderServiceName, args, "start", func() {}, opts)
}

func systemdUpgraderScript(planPath string) string {
	// 注意：这是 kardianos/service 的 systemd 模板文本，会被当作 text/template 解析，
	// 因此我们保留 {{.Path}} / {{.Arguments}} 等占位符，只把 planPath 写死进 ConditionPathExists。
	return fmt.Sprintf(`[Unit]
Description=frp-panel upgrader (oneshot)
ConditionFileIsExecutable={{.Path|cmdEscape}}
ConditionPathExists=%s

[Service]
Type=oneshot
ExecStart={{.Path|cmdEscape}}{{range .Arguments}} {{.|cmd}}{{end}}
{{if .WorkingDirectory}}WorkingDirectory={{.WorkingDirectory|cmdEscape}}{{end}}
Restart=no

[Install]
WantedBy=multi-user.target
`, planPath)
}

func mustExecutablePath() string {
	p, _ := os.Executable()
	if real, err := filepath.EvalSymlinks(p); err == nil && len(real) > 0 {
		p = real
	}
	if abs, err := filepath.Abs(p); err == nil {
		p = abs
	}
	return p
}

// RunWorker 执行升级计划（给隐藏命令 __upgrade-worker 调用）
func RunWorker(ctx context.Context, planPath string) error {
	if ctx == nil {
		ctx = context.Background()
	}

	// upgrader service 可能在 boot 或无 plan 时被启动：此时直接退出，保证不会误停 frpp
	if _, err := os.Stat(planPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	p, err := readPlan(planPath)
	if err != nil {
		return err
	}
	opt := p.Options

	target, err := resolveTargetPath(opt.TargetPath)
	if err != nil {
		return err
	}
	opt.TargetPath = target

	// worker 假设 staged 固定为 targetPath+".new"
	staged := stagePathForTarget(opt.TargetPath)
	if err := verifyBinary(staged); err != nil {
		_ = writeStatus(opt.WorkDir, false, err.Error())
		return err
	}

	if opt.RestartService && len(strings.TrimSpace(opt.ServiceName)) > 0 {
		// stop（这一步会导致 remoteshell 断开，但 worker 在独立 service 里执行，不会被一起杀掉）
		if err := utils.ControlSystemService(opt.ServiceName, opt.ServiceArgs, "stop", func() {}); err != nil {
			_ = writeStatus(opt.WorkDir, false, err.Error())
			return err
		}
	}

	// replace：必须避免打开/truncate target（会触发 ETXTBSY），优先用 rename（staged 与 target 同目录）
	if err := replaceStagedToTarget(ctx, staged, opt.TargetPath, opt.Backup); err != nil {
		_ = writeStatus(opt.WorkDir, false, err.Error())
		return err
	}

	if opt.RestartService && len(strings.TrimSpace(opt.ServiceName)) > 0 {
		if err := utils.ControlSystemService(opt.ServiceName, opt.ServiceArgs, "start", func() {}); err != nil {
			_ = writeStatus(opt.WorkDir, false, err.Error())
			return err
		}
	}

	_ = os.Remove(planPath) // 完成后清理 plan，避免开机自启误触发
	_ = writeStatus(opt.WorkDir, true, "ok")
	return nil
}

func replaceStagedToTarget(ctx context.Context, staged, target string, backup bool) error {
	_ = ctx
	if len(staged) == 0 || len(target) == 0 {
		return fmt.Errorf("staged/target is empty")
	}

	if runtime.GOOS == "windows" {
		// Windows：需要在 service stop 后才能动目标文件，且 rename 覆盖通常不允许
		deadline := time.Now().Add(60 * time.Second)
		var backupPath string
		for {
			if time.Now().After(deadline) {
				return fmt.Errorf("timeout waiting target file unlock: %s", target)
			}

			if backup && len(backupPath) == 0 {
				backupPath = target + ".bak"
				_ = os.Remove(backupPath)
				_ = os.Rename(target, backupPath)
			} else if !backup {
				_ = os.Remove(target)
			}

			if err := os.Rename(staged, target); err != nil {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			return nil
		}
	}

	// Unix-like：rename 覆盖是原子操作，且不受正在运行的旧 binary 影响（不会触发 ETXTBSY）
	if backup {
		backupPath := target + ".bak"
		_ = os.Remove(backupPath)
		// 备份采用 rename，避免 open/truncate
		_ = os.Rename(target, backupPath)
	}
	if err := os.Rename(staged, target); err != nil {
		return err
	}
	return nil
}
