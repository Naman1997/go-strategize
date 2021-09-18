package services

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Exists(path string, homedir string) (bool, error) {
	path = HomeFix(path, homedir)
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	return false, err
}

func HomeFix(path string, homedir string) string {
	if strings.Contains(path, "~/") {
		path = filepath.Join(homedir, path[2:])
	}
	return path
}

func Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
