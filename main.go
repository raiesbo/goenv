package goenv

import (
	"os"
	"path/filepath"
	"strings"
)

// getPaths receives an array of directory entries and transforms it into a array of file paths concatenating the
// directory name with a base path.
func getPaths(newDirs []os.DirEntry, basePath string) []string {
	var paths []string

	for _, dir := range newDirs {
		path := filepath.Join(basePath, dir.Name())
		paths = append(paths, path)
	}

	return paths
}

// loadVarsFromFile receives a path to an .env file, parses it and loads all the variables.
func loadVarsFromFile(path string) error {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	for _, line := range strings.Split(string(fileData), "\n") {
		lineContent := strings.Split(line, "=")
		if len(lineContent) == 2 {
			if err := os.Setenv(lineContent[0], lineContent[1]); err != nil {
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

	dirsQueue = getPaths(dirs, ".")

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
			dirsQueue = append(dirsQueue, getPaths(children, path)...)
		} else if file.Name() == ".env" {
			return loadVarsFromFile(path)
		}

		dirsQueue = dirsQueue[1:]
	}

	return nil
}
