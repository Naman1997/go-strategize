package services

import "os"

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil && len(path) > 0 {
		return true, nil
	}
	return false, err
}
