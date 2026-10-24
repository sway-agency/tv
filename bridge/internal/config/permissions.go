package config

import (
	"fmt"
	"os"
	"path/filepath"
)

func CheckPrivate(path string) error {
	for _, target := range []string{path, filepath.Dir(path)} {
		info, err := os.Stat(target)
		if err != nil {
			return err
		}
		if info.Mode().Perm()&0o077 != 0 {
			return fmt.Errorf("%s is readable by other users, run: chmod go-rwx %s", target, target)
		}
	}
	return nil
}
