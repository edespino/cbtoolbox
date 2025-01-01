package coreinfo

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// MockableLookPath defines a function signature for LookPath to allow mocking.
type MockableLookPath func(string) (string, error)

var originalCheckPrerequisites func() error

func TestMain(m *testing.M) {
	// Save the original implementation
	originalCheckPrerequisites = checkPrerequisites

	// Run tests
	code := m.Run()

	// Restore the original implementation
	checkPrerequisites = originalCheckPrerequisites

	os.Exit(code)
}

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
	tempDir := t.TempDir()

	// Create mock core files with ELF magic number
	coreFile1 := filepath.Join(tempDir, "core.1234")
	coreFile2 := filepath.Join(tempDir, "core")

	elfMagic := []byte("\x7fELF") // ELF magic number

	err := os.WriteFile(coreFile1, elfMagic, 0644)
	if err != nil {
		t.Fatalf("Failed to write mock core file1: %v", err)
	}

	err = os.WriteFile(coreFile2, elfMagic, 0644)
	if err != nil {
		t.Fatalf("Failed to write mock core file2: %v", err)
	}

	// Create a non-core file
	invalidFile := filepath.Join(tempDir, "invalid.txt")
	err = os.WriteFile(invalidFile, []byte("This is not a core file"), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid file: %v", err)
	}

	// Validate core files
	files, err := validateCoreFiles([]string{tempDir})
	if err != nil {
		t.Errorf("Unexpected error during validation: %v", err)
	}

	// Check expected results
	if len(files) != 2 {
		t.Errorf("Expected 2 core files, got %d", len(files))
	}
}

func TestCoreInfoVerboseOutput(t *testing.T) {
	// Mock checkPrerequisites to always succeed
	checkPrerequisites = func() error {
		return nil
	}

	// Check if gdb is available in PATH
	if _, err := exec.LookPath("gdb"); err != nil {
		t.Skip("gdb not found in PATH, skipping test on macOS")
	}

	// Attempt to find a real core file in /var/crash
	corePattern := "/var/crash/core-crash-*"
	matches, err := filepath.Glob(corePattern)
	var coreFiles []string

	if err == nil && len(matches) > 0 {
		// Use the real core files found
		t.Logf("Using real core file(s) for test: %v", matches)
		coreFiles = matches
	} else {
		// Fall back to creating mock core files
		t.Log("No real core files found, falling back to mock core files.")

		tempDir := t.TempDir()

		// Create mock core files with ELF magic number
		coreFile1 := filepath.Join(tempDir, "core.1234")
		coreFile2 := filepath.Join(tempDir, "core")

		elfMagic := []byte("\x7fELF") // ELF magic number

		err := os.WriteFile(coreFile1, elfMagic, 0644)
		if err != nil {
			t.Fatalf("Failed to write mock core file1: %v", err)
		}

		err = os.WriteFile(coreFile2, elfMagic, 0644)
		if err != nil {
			t.Fatalf("Failed to write mock core file2: %v", err)
		}

		coreFiles = []string{coreFile1, coreFile2}
	}

	// Capture verbose output
	verbose = true
	output := captureOutput(func() {
		err := RunCoreInfo(nil, coreFiles)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	fmt.Printf("Captured GDB Output:\n%s\n", output)

	// Validate verbose output
	for _, coreFile := range coreFiles {
		if !strings.Contains(output, fmt.Sprintf("Validating file: %s -> Valid core file", coreFile)) {
			t.Errorf("Expected verbose output for coreFile %s, got:\n%s", coreFile, output)
		}
	}

	// Validate summary output
	if len(coreFiles) > 1 {
		if !strings.Contains(output, fmt.Sprintf("Validated core files: [%s %s]", coreFiles[0], coreFiles[1])) &&
			!strings.Contains(output, fmt.Sprintf("Validated core files: [%s %s]", coreFiles[1], coreFiles[0])) {
			t.Errorf("Expected summary output to contain core files in any order, got:\n%s", output)
		}
	} else if len(coreFiles) == 1 {
		if !strings.Contains(output, fmt.Sprintf("Validated core files: [%s]", coreFiles[0])) {
			t.Errorf("Expected summary output to contain single core file, got:\n%s", output)
		}
	}
}

func captureOutput(f func()) string {
	r, w, _ := os.Pipe()
	stdOut := os.Stdout
	os.Stdout = w

	f()

	w.Close()
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to capture output: %v\n", err)
	}
	os.Stdout = stdOut

	return buf.String()
}
