//go:build !windows

package upgrade

import (
	"debug/elf"
	"debug/macho"
	"fmt"
	"runtime"
)

func verifyBinary(path string) error {
	if err := verifyFileNonEmpty(path); err != nil {
		return err
	}

	switch runtime.GOOS {
	case "linux":
		f, err := elf.Open(path)
		if err != nil {
			return fmt.Errorf("invalid ELF: %w", err)
		}
		defer f.Close()
		// 简单架构匹配（尽量不误伤）
		switch runtime.GOARCH {
		case "amd64":
			if f.FileHeader.Machine != elf.EM_X86_64 {
				return fmt.Errorf("arch mismatch: want amd64, got %v", f.FileHeader.Machine)
			}
		case "arm64":
			if f.FileHeader.Machine != elf.EM_AARCH64 {
				return fmt.Errorf("arch mismatch: want arm64, got %v", f.FileHeader.Machine)
			}
		case "arm":
			if f.FileHeader.Machine != elf.EM_ARM {
				return fmt.Errorf("arch mismatch: want arm, got %v", f.FileHeader.Machine)
			}
		}
		return nil
	case "darwin":
		f, err := macho.Open(path)
		if err != nil {
			return fmt.Errorf("invalid Mach-O: %w", err)
		}
		defer f.Close()
		switch runtime.GOARCH {
		case "amd64":
			if f.Cpu != macho.CpuAmd64 {
				return fmt.Errorf("arch mismatch: want amd64, got %v", f.Cpu)
			}
		case "arm64":
			if f.Cpu != macho.CpuArm64 {
				return fmt.Errorf("arch mismatch: want arm64, got %v", f.Cpu)
			}
		}
		return nil
	default:
		// 其他 Unix-like：仅做非空校验
		return nil
	}
}


