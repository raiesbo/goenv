package goenv

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// setupTempDir creates a temporary directory for testing
func setupTempDir(t *testing.T) (string, func()) {
	t.Helper()
	baseDir, err := os.MkdirTemp("", "goenv_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	cleanup := func() {
		os.RemoveAll(baseDir)
	}

	return baseDir, cleanup
}

// createEnvFile creates a .env file with the given content
func createEnvFile(t *testing.T, dir, filename, content string) string {
	t.Helper()
	envPath := filepath.Join(dir, filename)
	if err := os.WriteFile(envPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create %s file: %v", filename, err)
	}
	return envPath
}

// Test_Load tests the basic loading functionality with backward compatibility
func Test_Load(t *testing.T) {
	baseDir, cleanup := setupTempDir(t)
	defer cleanup()

	// Change to test directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(baseDir)

	subfolders := []string{"folder1", "folder2", "folder3"}
	var targetFolder string

	for _, folder := range subfolders {
		path := filepath.Join(baseDir, folder)
		if err := os.Mkdir(path, 0755); err != nil {
			t.Fatalf("Failed to create subfolder %s: %v", folder, err)
		}
		targetFolder = path
	}

	envContent := "TESTKEY=TEST_VALUE\nTEST_KEY_2=TEST_VALUE_2"
	createEnvFile(t, targetFolder, ".env", envContent)

	// Clear any existing env vars
	os.Unsetenv("TESTKEY")
	os.Unsetenv("TEST_KEY_2")

	if err := Load(); err != nil {
		t.Fatalf("Failed to find the .env file: %v", err)
	}

	if val := os.Getenv("TESTKEY"); val != "TEST_VALUE" {
		t.Errorf("Expected TESTKEY=TEST_VALUE, got %s", val)
	}

	if val := os.Getenv("TEST_KEY_2"); val != "TEST_VALUE_2" {
		t.Errorf("Expected TEST_KEY_2=TEST_VALUE_2, got %s", val)
	}
}

// Test_LoadWithConfig tests the new configuration-based loading
func Test_LoadWithConfig(t *testing.T) {
	baseDir, cleanup := setupTempDir(t)
	defer cleanup()

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(baseDir)

	// Create multiple env files
	createEnvFile(t, baseDir, ".env", "BASE_VAR=base_value")
	createEnvFile(t, baseDir, ".env.local", "LOCAL_VAR=local_value")

	config := &Config{
		EnvFiles:    []string{".env", ".env.local"},
		MaxDepth:    5,
		StopOnFirst: false,
	}

	// Clear existing vars
	os.Unsetenv("BASE_VAR")
	os.Unsetenv("LOCAL_VAR")

	if err := LoadWithConfig(config); err != nil {
		t.Fatalf("LoadWithConfig failed: %v", err)
	}

	if val := os.Getenv("BASE_VAR"); val != "base_value" {
		t.Errorf("Expected BASE_VAR=base_value, got %s", val)
	}

	if val := os.Getenv("LOCAL_VAR"); val != "local_value" {
		t.Errorf("Expected LOCAL_VAR=local_value, got %s", val)
	}
}

// Test_LoadFile tests loading a specific file
func Test_LoadFile(t *testing.T) {
	baseDir, cleanup := setupTempDir(t)
	defer cleanup()

	envContent := "SPECIFIC_VAR=specific_value"
	envPath := createEnvFile(t, baseDir, "custom.env", envContent)

	os.Unsetenv("SPECIFIC_VAR")

	if err := LoadFile(envPath); err != nil {
		t.Fatalf("LoadFile failed: %v", err)
	}

	if val := os.Getenv("SPECIFIC_VAR"); val != "specific_value" {
		t.Errorf("Expected SPECIFIC_VAR=specific_value, got %s", val)
	}
}

// Test_QuotedValues tests handling of quoted values
func Test_QuotedValues(t *testing.T) {
	baseDir, cleanup := setupTempDir(t)
	defer cleanup()

	envContent := `QUOTED_SINGLE='single quoted value'
QUOTED_DOUBLE="double quoted value"
UNQUOTED=unquoted value
EMPTY_QUOTED=""
SPACES_QUOTED="  spaces  "`

	envPath := createEnvFile(t, baseDir, ".env", envContent)

	// Clear vars
	vars := []string{"QUOTED_SINGLE", "QUOTED_DOUBLE", "UNQUOTED", "EMPTY_QUOTED", "SPACES_QUOTED"}
	for _, v := range vars {
		os.Unsetenv(v)
	}

	if err := LoadFile(envPath); err != nil {
		t.Fatalf("LoadFile failed: %v", err)
	}

	tests := []struct {
		key      string
		expected string
	}{
		{"QUOTED_SINGLE", "single quoted value"},
		{"QUOTED_DOUBLE", "double quoted value"},
		{"UNQUOTED", "unquoted value"},
		{"EMPTY_QUOTED", ""},
		{"SPACES_QUOTED", "  spaces  "},
	}

	for _, test := range tests {
		if val := os.Getenv(test.key); val != test.expected {
			t.Errorf("Expected %s=%q, got %q", test.key, test.expected, val)
		}
	}
}

// Test_InvalidFormat tests error handling for invalid .env format
func Test_InvalidFormat(t *testing.T) {
	baseDir, cleanup := setupTempDir(t)
	defer cleanup()

	envContent := "INVALID_LINE_NO_EQUALS"
	envPath := createEnvFile(t, baseDir, ".env", envContent)

	err := LoadFile(envPath)
	if err == nil {
		t.Fatal("Expected error for invalid format, got nil")
	}

	if !strings.Contains(err.Error(), "invalid format") {
		t.Errorf("Expected 'invalid format' error, got: %v", err)
	}
}

// Test_CycleDetection tests that the loader handles symbolic links properly
func Test_CycleDetection(t *testing.T) {
	baseDir, cleanup := setupTempDir(t)
	defer cleanup()

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(baseDir)

	// Create a deep directory structure
	deepDir := filepath.Join(baseDir, "level1", "level2", "level3", "level4", "level5", "level6")
	if err := os.MkdirAll(deepDir, 0755); err != nil {
		t.Fatalf("Failed to create deep directory: %v", err)
	}

	createEnvFile(t, deepDir, ".env", "DEEP_VAR=deep_value")

	config := &Config{
		EnvFiles:    []string{".env"},
		MaxDepth:    3, // Limit depth
		StopOnFirst: true,
	}

	os.Unsetenv("DEEP_VAR")

	// Should not find the file due to depth limit
	if err := LoadWithConfig(config); err != nil {
		t.Fatalf("LoadWithConfig failed: %v", err)
	}

	// Variable should not be set due to depth limit
	if val := os.Getenv("DEEP_VAR"); val != "" {
		t.Errorf("Expected DEEP_VAR to be empty due to depth limit, got %s", val)
	}

	// Now test with sufficient depth
	config.MaxDepth = 10
	if err := LoadWithConfig(config); err != nil {
		t.Fatalf("LoadWithConfig failed: %v", err)
	}

	if val := os.Getenv("DEEP_VAR"); val != "deep_value" {
		t.Errorf("Expected DEEP_VAR=deep_value, got %s", val)
	}
}

// Test_GetString verifies that GetString correctly retrieves environment variables and falls back when necessary.
func Test_GetString(t *testing.T) {
	key := "TEST_STRING"
	value := "hello"
	fallback := "default"

	os.Setenv(key, value)
	defer os.Unsetenv(key)

	if got := GetString(key, fallback); got != value {
		t.Errorf("GetString(%q, %q) = %q; want %q", key, fallback, got, value)
	}

	os.Unsetenv(key)
	if got := GetString(key, fallback); got != fallback {
		t.Errorf("GetString(%q, %q) = %q; want %q", key, fallback, got, fallback)
	}
}

// Test_GetInt verifies that GetInt correctly retrieves integer environment variables and falls back when necessary.
func Test_GetInt(t *testing.T) {
	key := "TEST_INT"
	value := "42"
	fallback := 10

	os.Setenv(key, value)
	defer os.Unsetenv(key)

	if got := GetInt(key, fallback); got != 42 {
		t.Errorf("GetInt(%q, %d) = %d; want %d", key, fallback, got, 42)
	}

	os.Setenv(key, "invalid")
	if got := GetInt(key, fallback); got != fallback {
		t.Errorf("GetInt(%q, %d) with invalid value = %d; want %d", key, fallback, got, fallback)
	}

	os.Unsetenv(key)
	if got := GetInt(key, fallback); got != fallback {
		t.Errorf("GetInt(%q, %d) = %d; want %d", key, fallback, got, fallback)
	}
}

// Test_GetBool tests the new GetBool function
func Test_GetBool(t *testing.T) {
	key := "TEST_BOOL"
	fallback := false

	tests := []struct {
		value    string
		expected bool
	}{
		{"true", true},
		{"True", true},
		{"TRUE", true},
		{"1", true},
		{"false", false},
		{"False", false},
		{"FALSE", false},
		{"0", false},
		{"invalid", fallback}, // Should fall back on invalid value
	}

	for _, test := range tests {
		os.Setenv(key, test.value)
		if got := GetBool(key, fallback); got != test.expected {
			t.Errorf("GetBool(%q, %t) with value %q = %t; want %t", key, fallback, test.value, got, test.expected)
		}
	}

	// Test fallback when key doesn't exist
	os.Unsetenv(key)
	if got := GetBool(key, true); got != true {
		t.Errorf("GetBool(%q, %t) with missing key = %t; want %t", key, true, got, true)
	}
}

// Test_GetFloat tests the new GetFloat function
func Test_GetFloat(t *testing.T) {
	key := "TEST_FLOAT"
	fallback := 3.14

	os.Setenv(key, "2.718")
	defer os.Unsetenv(key)

	if got := GetFloat(key, fallback); got != 2.718 {
		t.Errorf("GetFloat(%q, %f) = %f; want %f", key, fallback, got, 2.718)
	}

	os.Setenv(key, "invalid")
	if got := GetFloat(key, fallback); got != fallback {
		t.Errorf("GetFloat(%q, %f) with invalid value = %f; want %f", key, fallback, got, fallback)
	}

	os.Unsetenv(key)
	if got := GetFloat(key, fallback); got != fallback {
		t.Errorf("GetFloat(%q, %f) = %f; want %f", key, fallback, got, fallback)
	}
}

// Test_MustGetString tests the new MustGetString function
func Test_MustGetString(t *testing.T) {
	key := "TEST_MUST_STRING"
	value := "required_value"

	os.Setenv(key, value)
	defer os.Unsetenv(key)

	if got := MustGetString(key); got != value {
		t.Errorf("MustGetString(%q) = %q; want %q", key, got, value)
	}

	os.Unsetenv(key)

	// Test panic behavior
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustGetString should have panicked for missing key")
		}
	}()

	MustGetString(key) // Should panic
}

// Test_CommentsAndEmptyLines tests handling of comments and empty lines
func Test_CommentsAndEmptyLines(t *testing.T) {
	baseDir, cleanup := setupTempDir(t)
	defer cleanup()

	envContent := `# This is a comment
VAR1=value1

# Another comment
VAR2=value2

# Empty line above and below

VAR3=value3
# Inline comment - this line should be ignored completely`

	envPath := createEnvFile(t, baseDir, ".env", envContent)

	// Clear vars
	vars := []string{"VAR1", "VAR2", "VAR3"}
	for _, v := range vars {
		os.Unsetenv(v)
	}

	if err := LoadFile(envPath); err != nil {
		t.Fatalf("LoadFile failed: %v", err)
	}

	expected := map[string]string{
		"VAR1": "value1",
		"VAR2": "value2",
		"VAR3": "value3",
	}

	for key, expectedVal := range expected {
		if val := os.Getenv(key); val != expectedVal {
			t.Errorf("Expected %s=%s, got %s", key, expectedVal, val)
		}
	}
}

// Test_StopOnFirst tests the StopOnFirst configuration option
func Test_StopOnFirst(t *testing.T) {
	baseDir, cleanup := setupTempDir(t)
	defer cleanup()

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(baseDir)

	// Create nested structure with multiple .env files
	subDir := filepath.Join(baseDir, "subdir")
	os.Mkdir(subDir, 0755)

	createEnvFile(t, baseDir, ".env", "ROOT_VAR=root_value")
	createEnvFile(t, subDir, ".env", "SUB_VAR=sub_value")

	// Test StopOnFirst = true (default)
	config := &Config{
		EnvFiles:    []string{".env"},
		StopOnFirst: true,
	}

	os.Unsetenv("ROOT_VAR")
	os.Unsetenv("SUB_VAR")

	if err := LoadWithConfig(config); err != nil {
		t.Fatalf("LoadWithConfig failed: %v", err)
	}

	// Should only load the first .env found
	if val := os.Getenv("ROOT_VAR"); val != "root_value" {
		t.Errorf("Expected ROOT_VAR=root_value, got %s", val)
	}

	// SUB_VAR should not be loaded due to StopOnFirst
	if val := os.Getenv("SUB_VAR"); val != "" {
		t.Errorf("Expected SUB_VAR to be empty due to StopOnFirst, got %s", val)
	}
}

// Test_InvalidKeyNames tests validation of environment variable names
func Test_InvalidKeyNames(t *testing.T) {
	baseDir, cleanup := setupTempDir(t)
	defer cleanup()

	envContent := `VALID_KEY=valid_value
123INVALID=invalid_start_with_number
INVALID-DASH=invalid_dash
=no_key`

	envPath := createEnvFile(t, baseDir, ".env", envContent)

	err := LoadFile(envPath)
	if err == nil {
		t.Fatal("Expected error for invalid key names, got nil")
	}

	// Should contain information about invalid key
	if !strings.Contains(err.Error(), "invalid environment variable name") {
		t.Errorf("Expected 'invalid environment variable name' error, got: %v", err)
	}
}
