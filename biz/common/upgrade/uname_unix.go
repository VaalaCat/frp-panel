//go:build !windows

package upgrade

import (
	"strings"

	"golang.org/x/sys/unix"
)

func unameMachine() string {
	var u unix.Utsname
	if err := unix.Uname(&u); err != nil {
		return ""
	}
	return strings.TrimSpace(bytesToString(u.Machine[:]))
}

func bytesToString(ca []byte) string {
	b := make([]byte, 0, len(ca))
	for _, c := range ca {
		if c == 0 {
			break
		}
		b = append(b, c)
	}
	return string(b)
}


