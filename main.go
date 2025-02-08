package goenv

import (
	"fmt"
	"os"
	"strings"
)

// appendPaths receives an array of directory entries and appends their name with full path to a predefined
// directories queue.
func appendPaths(dirsQueue *[]string, newDirs []os.DirEntry, basePath string) {
	for _, dir := range newDirs {
		filePath := fmt.Sprintf("%s/%s", basePath, dir.Name())
		*dirsQueue = append(*dirsQueue, filePath)
	}
}

// loadVarsFromFile receives a path to an .env file, parses it and loads all env variables.
func loadVarsFromFile(path string) error {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	for _, line := range strings.Split(string(fileData), "\n") {
		lineContent := strings.Split(line, "=")
		if len(lineContent) == 2 {
			err := os.Setenv(lineContent[0], lineContent[1])
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Load reads recursively all the directories of a project until it finds a .env file. Once the .env is found, reads
// the file and loads the values as OS ENV values.
func Load() error {
	var dirsQueue []string

	dirs, err := os.ReadDir("./")
	if err != nil {
		return err
	}

	appendPaths(&dirsQueue, dirs, ".")

	for len(dirsQueue) > 0 {
		path := dirsQueue[0]

		file, err := os.Stat(path)
		if err != nil {
			return err
		}

		if file.IsDir() {
			children, err := os.ReadDir(path)
			if err != nil {
				return err
			}
			appendPaths(&dirsQueue, children, path)
		} else if file.Name() == ".env" {
			return loadVarsFromFile(path)
		}

		dirsQueue = dirsQueue[1:]
	}

	return nil
}
