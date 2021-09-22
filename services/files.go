package services

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

/*
Exists checks if path provided exists or not.
Exits with an error if not.
*/
func Exists(path string, homedir string) bool {
	_, err := os.Stat(path)
	if err != nil {
		ColorPrint(ERROR, err.Error())
	}
	return true
}

//ExistsFolderNoErr checks if path provided is a folder or not.
func ExistsFolderNoErr(path string, homedir string) bool {
	f, err := os.Stat(path)
	if err == nil && f.IsDir() {
		return true
	}
	return false
}

//HomeFix resolves home-relative path to absolute path
func HomeFix(path string, homedir string) string {
	if strings.Contains(path, "~/") {
		path = filepath.Join(homedir, path[2:])
	}
	return path
}

//Copy copies data from src to dst path
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

/*
Validate checks if a relative path exists or not.
Returns absolute path if its relative to home dir.
Exits with status 1 if path not found
*/
func Validate(path string, homedir string) string {
	path = HomeFix(path, homedir)
	_ = Exists(path, homedir)
	return path
}

/*
ReadFiles walks through all filepaths in a given dir.
Returns a list of filepaths if successful
Exists with an error otherwise
*/
func ReadFiles(searchDir string) ([]string, error) {
	fileList := make([]string, 0)
	e := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return err
	})

	if e != nil {
		ColorPrint(ERROR, e.Error())
	}

	return fileList, nil
}
