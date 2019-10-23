// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// FileExists returns true if filename exists
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// DirExists returns true if dirname exists
func DirExists(dirname string) bool {
	info, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// GetFilesDB returns files with extension the extension passed from the lists
func GetFilesDB(extension string, listFiles []string, listDirs []string) ([]string, error) {
	dbFiles := make([]string, 0)
	for _, dir := range listDirs {
		if !DirExists(dir) {
			return dbFiles, fmt.Errorf("directory '%s' doesn't exists", dir)
		}
		files, err := getFilesSuffix(dir, fmt.Sprintf(".%s", extension))
		if err != nil {
			return nil, fmt.Errorf("getting files from '%s': %v", dir, err)
		}
		for _, file := range files {
			dbFiles = append(dbFiles, filepath.Clean(dir+string(filepath.Separator)+file))
		}
	}
	for _, file := range listFiles {
		if !FileExists(file) {
			return dbFiles, fmt.Errorf("file '%s' doesn't exists", file)
		}
		if !strings.HasSuffix(file, fmt.Sprintf(".%s", extension)) {
			return dbFiles, fmt.Errorf("file '%s' without extension '%s'", file, extension)
		}
		dbFiles = append(dbFiles, filepath.Clean(file))
	}
	return dbFiles, nil
}

func getFilesSuffix(directory string, suffix string) ([]string, error) {
	readf, err := ioutil.ReadDir(directory)
	if err != nil {
		return []string{}, err
	}
	files := make([]string, 0, len(readf))
	for _, f := range readf {
		if !f.IsDir() {
			if strings.HasSuffix(f.Name(), suffix) {
				files = append(files, f.Name())
			}
		}
	}
	return files, nil
}
