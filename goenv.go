package goenv

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	EnvFile = ".env"
)

// getPaths converts an array of directory entries into an array of file paths by concatenating each
// entry's name with a base path.
func getPaths(newDirs []os.DirEntry, basePath string) []string {
	paths := make([]string, len(newDirs))

	for i, dir := range newDirs {
		paths[i] = filepath.Join(basePath, dir.Name())
	}

	return paths
}

// loadVarsFromFile parses an .env file at the given path and loads its variables into the environment.
func loadVarsFromFile(path string) error {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	for _, line := range strings.Split(string(fileData), "\n") {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, "=")
		if len(parts) == 2 {
			key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
			if err := os.Setenv(key, value); err != nil {
				return err
			}
		}
	}

	return nil
}

// Load recursively scans all directories of a project until it finds a .env file. Once found, it reads
// the file and loads its values as environment variables.
func Load() error {
	dirsQueue := []string{"./"}

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
		} else if file.Name() == EnvFile {
			return loadVarsFromFile(path)
		}

		dirsQueue = dirsQueue[1:]
	}

	return nil
}

// GetString returns the value of an environment variable if it exists; otherwise, it returns the fallback value.
func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return val
}

// GetInt returns the integer value of an environment variable if it exists; otherwise, it returns the fallback value.
func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return intVal
}
