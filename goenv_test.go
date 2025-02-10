package goenv

import (
	"os"
	"path/filepath"
	"testing"
)

// Test_Load tests the creation of directories and an .env file.
func Test_Load(t *testing.T) {
	baseDir, err := os.MkdirTemp("./", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(baseDir)

	subfolders := []string{"folder1", "folder2", "folder3"}
	var targetFolder string

	for _, folder := range subfolders {
		path := filepath.Join(baseDir, folder)
		if err := os.Mkdir(path, 0755); err != nil {
			t.Fatalf("Failed to create subfolder %s: %v", folder, err)
		}
		targetFolder = path
	}

	envFilePath := filepath.Join(targetFolder, ".env")
	envContent := "TEST_KEY=TEST_VALUE\nTEST_KEY_2=TEST_VALUE_2"

	if err := os.WriteFile(envFilePath, []byte(envContent), 0644); err != nil {
		t.Fatalf("Failed to create .env file: %v", err)
	}

	if err := Load(); err != nil {
		t.Fatalf("Failed to find the .env file: %v", err)
	}

	if val := os.Getenv("TEST_KEY"); val != "TEST_VALUE" {
		t.Fatalf("Failed to load env variables: %v", err)
	}

	if val := os.Getenv("TEST_KEY_2"); val != "TEST_VALUE_2" {
		t.Fatalf("Failed to load env variables: %v", err)
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
