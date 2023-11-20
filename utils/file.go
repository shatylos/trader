package utils

import "os"

func IsDir(filename string) (bool, error) {
	fd, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	fm := fd.Mode()
	return fm.IsDir(), nil
}

func MkDir(filename string, mode os.FileMode) error {
	return os.MkdirAll(filename, mode)
}
