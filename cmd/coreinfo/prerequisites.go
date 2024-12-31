package coreinfo

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// checkPrerequisites verifies that all necessary tools and configurations are available.
func checkPrerequisites() error {
	if err := checkGDBAvailability(); err != nil {
		return err
	}
	// Add more prerequisite checks here if needed
	return nil
}

// checkGDBAvailability checks if the gdb command is available in the system's PATH.
func checkGDBAvailability() error {
	_, err := exec.LookPath("gdb")
	if err != nil {
		return fmt.Errorf("gdb is not installed or not available in PATH")
	}
	return nil
}

// validateCoreFiles ensures that the provided arguments contain valid core files or directories.
func validateCoreFiles(args []string) ([]string, error) {
	if len(args) == 0 {
		return nil, errors.New("no core files or directories provided")
	}

	var coreFiles []string
	for _, arg := range args {
		info, err := os.Stat(arg)
		if err != nil {
			return nil, fmt.Errorf("error accessing path '%s': %v", arg, err)
		}

		if info.IsDir() {
			// Search for core files in the directory
			files, err := filepath.Glob(filepath.Join(arg, "core*"))
			if err != nil {
				return nil, fmt.Errorf("error scanning directory '%s': %v", arg, err)
			}
			if len(files) == 0 {
				return nil, fmt.Errorf("no core files found in directory '%s'", arg)
			}
			coreFiles = append(coreFiles, files...)
		} else if filepath.Base(arg) == "core" || strings.HasPrefix(filepath.Base(arg), "core.") {
			// Single core file
			coreFiles = append(coreFiles, arg)
		} else {
			return nil, fmt.Errorf("invalid core file or directory: '%s'", arg)
		}
	}

	if len(coreFiles) == 0 {
		return nil, errors.New("no valid core files provided")
	}
	return coreFiles, nil
}
