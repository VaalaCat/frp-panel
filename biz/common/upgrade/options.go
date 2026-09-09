package upgrade

// Options 升级参数（与 CLI flags 对应）
type Options struct {
	// Version 目标版本，默认 latest（仅当 DownloadURL 为空时生效）
	Version string

	// DownloadURL 自定义下载地址（最高优先级）。若非空，则直接下载该 URL
	DownloadURL string

	// GithubProxy GitHub 代理前缀，例如 https://ghfast.top/（仅当 UseGithubProxy=true 且 DownloadURL 为空时生效）
	GithubProxy string

	// UseGithubProxy 是否启用 GithubProxy
	UseGithubProxy bool

	// HTTPProxy 下载用的 http/https 代理（透传给 req/v3）
	HTTPProxy string

	// TargetPath 要覆盖的可执行文件路径，默认当前运行的可执行文件路径（会尝试解析 symlink）
	TargetPath string

	// Backup 覆盖前是否备份旧文件（.bak）
	Backup bool

	// ServiceName 需要控制的服务名（为空则不做服务控制）
	ServiceName string

	// RestartService 是否在替换成功后重启服务（会导致服务短暂中断）
	RestartService bool

	// WorkDir 升级临时目录（plan/lock/download cache 等）
	WorkDir string

	// ServiceArgs 透传给 utils.ControlSystemService（参考 cmd/frpp/shared/cmd.go）
	ServiceArgs []string
}


