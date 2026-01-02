package upgrade

import (
	"fmt"
	"io"
	"os"
)

func copyFile(src, dst string, perm os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

func backupExisting(targetPath string) (string, error) {
	if _, err := os.Stat(targetPath); err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("stat existing binary failed: %w", err)
	}
	backupPath := targetPath + ".bak"
	_ = os.Remove(backupPath)
	if err := copyFile(targetPath, backupPath, 0755); err != nil {
		return "", fmt.Errorf("backup existing binary failed: %w", err)
	}
	return backupPath, nil
}

func replaceFile(src, dst string) error {
	if err := os.Rename(src, dst); err == nil {
		return nil
	}
	// rename 失败则退化为 copy（跨文件系统等）
	if err := copyFile(src, dst, 0755); err != nil {
		return err
	}
	return nil
}
