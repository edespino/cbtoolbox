package coreinfo

import (
	"bytes"
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

// isCoreFile determines if a file is a core dump using the `file` command.
func isCoreFile(filePath string) (bool, error) {
	cmd := exec.Command("file", filePath)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return false, fmt.Errorf("failed to run 'file' command: %w", err)
	}

	output := out.String()
	return strings.Contains(output, "core file"), nil
}

// validateCoreFiles validates the input paths to determine if they are core files or directories containing core files.
func validateCoreFiles(args []string) ([]string, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("no core files or directories provided")
	}

	var coreFiles []string
	for _, arg := range args {
		info, err := os.Stat(arg)
		if err != nil {
			return nil, fmt.Errorf("error accessing path '%s': %v", arg, err)
		}

		if info.IsDir() {
			// Search for files in the directory
			files, err := filepath.Glob(filepath.Join(arg, "*"))
			if err != nil {
				return nil, fmt.Errorf("error scanning directory '%s': %v", arg, err)
			}
			for _, file := range files {
				if valid, err := isCoreFile(file); err == nil && valid {
					coreFiles = append(coreFiles, file)
				}
			}
		} else {
			// Validate single file
			if valid, err := isCoreFile(arg); err == nil && valid {
				coreFiles = append(coreFiles, arg)
			} else {
				return nil, fmt.Errorf("invalid core file: '%s'", arg)
			}
		}
	}

	if len(coreFiles) == 0 {
		return nil, fmt.Errorf("no valid core files provided")
	}
	return coreFiles, nil
}
