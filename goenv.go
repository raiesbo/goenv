package goenv

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	DefaultEnvFile = ".env"
	MaxDepth       = 10 // Prevent infinite recursion
)

type Config struct {
	EnvFiles    []string
	MaxDepth    int
	StopOnFirst bool
	Prefix      string
}

func DefaultConfig() *Config {
	return &Config{
		EnvFiles:    []string{".env"},
		MaxDepth:    MaxDepth,
		StopOnFirst: true,
		Prefix:      "",
	}
}

// LoadWithConfig loads environment variables with custom configuration
func LoadWithConfig(config *Config) error {
	if config == nil {
		config = DefaultConfig()
	}

	return loadFromDirectory(".", config, 0, make(map[string]bool))
}

// Load provides backward compatibility with default behavior
func Load() error {
	return LoadWithConfig(nil)
}

// loadFromDirectory recursively searches for .env files with proper error handling and cycle detection
func loadFromDirectory(dir string, config *Config, depth int, visited map[string]bool) error {
	if depth > config.MaxDepth {
		return nil
	}

	// Get absolute path to detect cycles
	absPath, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for %s: %w", dir, err)
	}

	if visited[absPath] {
		return nil // Skip already visited directories
	}
	visited[absPath] = true

	// Check for .env files in current directory
	for _, envFile := range config.EnvFiles {
		envPath := filepath.Join(dir, envFile)
		if fileExists(envPath) {
			if err := loadVarsFromFile(envPath); err != nil {
				return fmt.Errorf("failed to load %s: %w", envPath, err)
			}
			if config.StopOnFirst {
				return nil
			}
		}
	}

	// Read directory entries
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	// Recursively search subdirectories
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			subDir := filepath.Join(dir, entry.Name())
			if err := loadFromDirectory(subDir, config, depth+1, visited); err != nil {
				return err
			}
		}
	}

	return nil
}

// fileExists checks if a file exists and is not a directory
func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// loadVarsFromFile parses an .env file with improved error handling and format support
func loadVarsFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key=value pairs
		if err := parseAndSetEnvVar(line, path, lineNum); err != nil {
			return err
		}
	}

	return scanner.Err()
}

// parseAndSetEnvVar parses a single environment variable line
func parseAndSetEnvVar(line, filePath string, lineNum int) error {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid format in %s at line %d: %s", filePath, lineNum, line)
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	// Handle quoted values
	if len(value) >= 2 {
		if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
			(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
			value = value[1 : len(value)-1]
		}
	}

	// Validate key format
	if !isValidEnvKey(key) {
		return fmt.Errorf("invalid environment variable name in %s at line %d: %s", filePath, lineNum, key)
	}

	return os.Setenv(key, value)
}

// isValidEnvKey validates environment variable key format
func isValidEnvKey(key string) bool {
	if key == "" {
		return false
	}

	for i, r := range key {
		if i == 0 {
			if !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == '_') {
				return false
			}
		} else {
			if !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_') {
				return false
			}
		}
	}

	return true
}

// GetString returns the value of an environment variable with fallback
func GetString(key, fallback string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return fallback
}

// GetInt returns the integer value of an environment variable with fallback
func GetInt(key string, fallback int) int {
	val, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}

	if intVal, err := strconv.Atoi(val); err == nil {
		return intVal
	}

	return fallback
}

// GetBool returns the boolean value of an environment variable with fallback
func GetBool(key string, fallback bool) bool {
	val, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}

	if boolVal, err := strconv.ParseBool(val); err == nil {
		return boolVal
	}

	return fallback
}

// GetFloat returns the float64 value of an environment variable with fallback
func GetFloat(key string, fallback float64) float64 {
	val, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}

	if floatVal, err := strconv.ParseFloat(val, 64); err == nil {
		return floatVal
	}

	return fallback
}

// MustGetString returns the value of an environment variable or panics if not found
func MustGetString(key string) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		panic(fmt.Sprintf("required environment variable %s not found", key))
	}
	return val
}

// LoadFile loads a specific .env file
func LoadFile(path string) error {
	return loadVarsFromFile(path)
}

// Unload removes all environment variables loaded from .env files
// Note: This is a simplified implementation - tracking loaded vars would be better
func Unload(keys []string) error {
	for _, key := range keys {
		if err := os.Unsetenv(key); err != nil {
			return err
		}
	}
	return nil
}
