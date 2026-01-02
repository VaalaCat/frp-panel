package upgrade

import (
	"fmt"
	"os"
	"path/filepath"
)

func resolveTargetPath(targetPath string) (string, error) {
	if len(targetPath) == 0 {
		exePath, err := os.Executable()
		if err != nil {
			return "", fmt.Errorf("获取当前执行文件失败: %w", err)
		}
		// 优先解析符号链接，确保替换真实文件
		if realPath, err := filepath.EvalSymlinks(exePath); err == nil && len(realPath) > 0 {
			exePath = realPath
		}
		targetPath = exePath
	}
	if abs, err := filepath.Abs(targetPath); err == nil {
		targetPath = abs
	}
	return targetPath, nil
}

func stagePathForTarget(targetPath string) string {
	return targetPath + ".new"
}
