package upgrade

import (
	"fmt"
	"os"
)

func verifyFileNonEmpty(path string) error {
	st, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat failed: %w", err)
	}
	if st.Size() <= 0 {
		return fmt.Errorf("file is empty: %s", path)
	}
	return nil
}
