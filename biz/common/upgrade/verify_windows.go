//go:build windows

package upgrade

import (
	"debug/pe"
	"fmt"
	"runtime"
)

func verifyBinary(path string) error {
	if err := verifyFileNonEmpty(path); err != nil {
		return err
	}
	f, err := pe.Open(path)
	if err != nil {
		return fmt.Errorf("invalid PE: %w", err)
	}
	defer f.Close()

	switch runtime.GOARCH {
	case "amd64":
		if f.Machine != pe.IMAGE_FILE_MACHINE_AMD64 {
			return fmt.Errorf("arch mismatch: want amd64, got 0x%x", f.Machine)
		}
	case "arm64":
		if f.Machine != pe.IMAGE_FILE_MACHINE_ARM64 {
			return fmt.Errorf("arch mismatch: want arm64, got 0x%x", f.Machine)
		}
	}
	return nil
}
