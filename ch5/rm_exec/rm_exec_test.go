package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRemoveExec(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Create a test executable file
	execFile := filepath.Join(tempDir, "test_exec")
	err := os.WriteFile(execFile, []byte("test"), 0755)
	if err != nil {
		t.Fatalf("Failed to create test executable file: %v", err)
	}

	// Call the remove_executable function
	files, err := remove_executable(tempDir)
	if err != nil {
		t.Fatalf("remove_executable failed: %v", err)
	}

	// Check if the executable file was returned
	if len(files) == 0 || files[0] != execFile {
		t.Errorf("Expected executable file %s to be returned, got %v", execFile, files)
	}
}
