package coreinfo

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// MockableLookPath defines a function signature for LookPath to allow mocking.
type MockableLookPath func(string) (string, error)

// checkGDBAvailabilityMockable is a testable version of checkGDBAvailability.
func checkGDBAvailabilityMockable(lookPath MockableLookPath) error {
	_, err := lookPath("gdb")
	if err != nil {
		return fmt.Errorf("gdb is not installed or not available in PATH")
	}
	return nil
}

// TestCheckGDBAvailability tests the availability of gdb using mockable LookPath.
func TestCheckGDBAvailability(t *testing.T) {
	// Test case: gdb is available
	lookPathMock := func(file string) (string, error) {
		if file == "gdb" {
			return "/usr/bin/gdb", nil
		}
		return "", errors.New("not found")
	}

	err := checkGDBAvailabilityMockable(lookPathMock)
	if err != nil {
		t.Errorf("Expected gdb to be available, got error: %v", err)
	}

	// Test case: gdb is not available
	lookPathMock = func(file string) (string, error) {
		return "", errors.New("not found")
	}

	err = checkGDBAvailabilityMockable(lookPathMock)
	if err == nil {
		t.Errorf("Expected error for unavailable gdb, got nil")
	}
}

// TestValidateCoreFiles validates core file paths and directories.
func TestValidateCoreFiles(t *testing.T) {
	// Create temporary directories and files for testing
	tempDir := t.TempDir()

	// Mock core files
	coreFile1 := filepath.Join(tempDir, "core.1234")
	coreFile2 := filepath.Join(tempDir, "core")

	err := os.WriteFile(coreFile1, []byte("mock core data"), 0644)
	if err != nil {
		t.Fatalf("Failed to write mock core file1: %v", err)
	}

	err = os.WriteFile(coreFile2, []byte("mock core data"), 0644)
	if err != nil {
		t.Fatalf("Failed to write mock core file2: %v", err)
	}

	// Test case: Single core file
	files, err := validateCoreFiles([]string{coreFile1})
	if err != nil || len(files) != 1 {
		t.Errorf("Expected 1 core file, got %v, error: %v", len(files), err)
	}

	// Test case: Directory containing core files
	files, err = validateCoreFiles([]string{tempDir})
	if err != nil || len(files) != 2 {
		t.Errorf("Expected 2 core files, got %v, error: %v", len(files), err)
	}

	// Test case: Invalid file
	invalidFile := filepath.Join(tempDir, "invalid.txt")
	err = os.WriteFile(invalidFile, []byte("invalid"), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid file: %v", err)
	}

	_, err = validateCoreFiles([]string{invalidFile})
	if err == nil {
		t.Errorf("Expected error for invalid file, got nil")
	}

	// Test case: No arguments
	_, err = validateCoreFiles([]string{})
	if err == nil {
		t.Errorf("Expected error for no arguments, got nil")
	}
}
