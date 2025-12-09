package common

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/VaalaCat/frp-panel/utils"
	"github.com/VaalaCat/frp-panel/utils/logger"
)

// UpgradeOptions 定义自助升级所需的参数
type UpgradeOptions struct {
	// Version 指定要升级的版本，默认为 latest
	Version string
	// GithubProxy 形如 https://ghfast.top/ 的前缀，会直接拼在下载链接前
	GithubProxy string
	// HTTPProxy 传递给 req/v3，用于走 HTTP/HTTPS 代理
	HTTPProxy string
	// TargetPath 需要覆盖的可执行文件路径，默认为当前运行的 frp-panel 路径
	TargetPath string
	// Backup 覆盖前是否备份旧文件，默认 true
	Backup bool
	// StopService 升级前是否尝试停止 systemd 服务，避免二进制被占用
	StopService bool
	// ServiceName systemd 服务名，默认 frpp
	ServiceName string
	// UseGithubProxy 仅当显式开启时才使用 Github 代理
	UseGithubProxy bool
}

// UpgradeSelf 下载并替换当前可执行文件
func UpgradeSelf(ctx context.Context, opt UpgradeOptions) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if err := opt.fillDefaults(ctx); err != nil {
		return err
	}

	asset, err := detectAssetName()
	if err != nil {
		return err
	}

	var (
		backupPath       string
		serviceWasActive bool
	)

	if opt.StopService && len(opt.ServiceName) > 0 {
		serviceWasActive, err = stopServiceIfActive(ctx, opt.ServiceName)
		if err != nil {
			return err
		}
	}

	defer func() {
		// 失败回滚
		if err != nil && len(backupPath) > 0 {
			if rErr := restoreBackup(ctx, backupPath, opt.TargetPath); rErr != nil {
				logger.Logger(ctx).Warnf("failed to restore backup, please check manually: %v", rErr)
			}
		}
		// 按原状态决定是否重启
		if serviceWasActive {
			if startErr := controlService(ctx, "start", opt.ServiceName); startErr != nil {
				logger.Logger(ctx).Warnf("failed to start service after upgrade, please check manually: %v", startErr)
				if err == nil {
					err = startErr
				}
			}
		}
	}()

	downloadURL := fmt.Sprintf("https://github.com/VaalaCat/frp-panel/releases/download/%s/%s", opt.Version, asset)
	if opt.UseGithubProxy && len(opt.GithubProxy) > 0 {
		downloadURL = fmt.Sprintf("%s/%s", strings.TrimRight(opt.GithubProxy, "/"), downloadURL)
	}

	logger.Logger(ctx).Infof("start downloading version [%s], url: %s", opt.Version, downloadURL)
	tmpPath, err := utils.DownloadFile(ctx, downloadURL, opt.HTTPProxy)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	if err := os.Chmod(tmpPath, 0755); err != nil {
		logger.Logger(ctx).Warnf("set file permission failed: %v", err)
	}

	if err := utils.EnsureDirectoryExists(opt.TargetPath); err != nil {
		return fmt.Errorf("ensure target directory failed: %w", err)
	}

	if opt.Backup {
		if backupPath, err = backupExisting(ctx, opt.TargetPath); err != nil {
			return err
		}
	}

	if err := replaceFile(tmpPath, opt.TargetPath); err != nil {
		return fmt.Errorf("replace executable failed: %w", err)
	}

	logger.Logger(ctx).Infof("frp-panel upgraded successfully, path: %s", opt.TargetPath)
	return nil
}

func (opt *UpgradeOptions) fillDefaults(ctx context.Context) error {
	if len(opt.Version) == 0 {
		opt.Version = "latest"
	}
	if len(opt.TargetPath) == 0 {
		exePath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("获取当前执行文件失败: %w", err)
		}
		// 优先解析符号链接，确保替换真实文件
		if realPath, err := filepath.EvalSymlinks(exePath); err == nil && len(realPath) > 0 {
			exePath = realPath
		}
		opt.TargetPath = exePath
	}

	if absPath, err := filepath.Abs(opt.TargetPath); err == nil {
		opt.TargetPath = absPath
	}

	if opt.StopService && len(opt.ServiceName) == 0 {
		opt.ServiceName = "frpp"
	}

	// 允许用户显式传空字符串来禁用代理
	opt.GithubProxy = strings.TrimSpace(opt.GithubProxy)
	opt.HTTPProxy = strings.TrimSpace(opt.HTTPProxy)

	return nil
}

func detectAssetName() (string, error) {
	osName := runtime.GOOS
	arch := runtime.GOARCH

	unameArch := arch
	if runtime.GOOS != "windows" {
		if out, err := exec.Command("uname", "-m").Output(); err == nil {
			unameArch = strings.TrimSpace(string(out))
		}
	}

	switch osName {
	case "linux":
		switch unameArch {
		case "x86_64", "amd64":
			return "frp-panel-linux-amd64", nil
		case "aarch64", "arm64":
			return "frp-panel-linux-arm64", nil
		case "armv7l":
			return "frp-panel-linux-armv7l", nil
		case "armv6l":
			return "frp-panel-linux-armv6l", nil
		}
	case "darwin":
		switch unameArch {
		case "x86_64", "amd64":
			return "frp-panel-darwin-amd64", nil
		case "arm64":
			return "frp-panel-darwin-arm64", nil
		}
	case "windows":
		switch arch {
		case "amd64":
			return "frp-panel-windows-amd64.exe", nil
		case "arm64":
			return "frp-panel-windows-arm64.exe", nil
		}
	}

	return "", fmt.Errorf("暂不支持的系统/架构: %s %s", osName, unameArch)
}

func backupExisting(ctx context.Context, path string) (string, error) {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		return "", fmt.Errorf("stat existing binary failed: %w", err)
	}

	backupPath := path + ".bak"
	_ = os.Remove(backupPath)

	if err := copyFile(path, backupPath); err != nil {
		return "", fmt.Errorf("backup existing binary failed: %w", err)
	}

	logger.Logger(ctx).Infof("backup created at: %s", backupPath)
	return backupPath, nil
}

func replaceFile(src, dst string) error {
	if err := os.Rename(src, dst); err == nil {
		return nil
	}

	if err := copyFile(src, dst); err != nil {
		return err
	}

	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	return out.Sync()
}

func restoreBackup(ctx context.Context, backupPath, target string) error {
	if len(backupPath) == 0 {
		return nil
	}
	logger.Logger(ctx).Infof("attempt to restore from backup: %s -> %s", backupPath, target)
	return replaceFile(backupPath, target)
}

func controlService(ctx context.Context, action, serviceName string) error {
	cmd := exec.CommandContext(ctx, "systemctl", action, serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("systemctl %s %s failed: %w, output: %s", action, serviceName, err, string(output))
	}
	logger.Logger(ctx).Infof("systemctl %s %s success", action, serviceName)
	return nil
}

func stopServiceIfActive(ctx context.Context, serviceName string) (bool, error) {
	cmd := exec.CommandContext(ctx, "systemctl", "is-active", "--quiet", serviceName)
	if err := cmd.Run(); err != nil {
		// 非 active，无需停
		logger.Logger(ctx).Infof("service %s is not active, skip stop", serviceName)
		return false, nil
	}

	if err := controlService(ctx, "stop", serviceName); err != nil {
		return false, err
	}
	logger.Logger(ctx).Infof("service %s stopped, ready to upgrade", serviceName)
	return true, nil
}
