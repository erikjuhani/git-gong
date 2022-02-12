package fs

import (
	"errors"
	"os"
)

func EnsureDir(dirname string) error {
	if err := os.Mkdir(dirname, os.ModePerm); os.IsNotExist(err) {
		fileinfo, err := os.Stat(dirname)
		if err != nil {
			return err
		}

		if fileinfo.IsDir() {
			return nil
		}

		return errors.New("path is not a directory")
	}

	return nil
}
