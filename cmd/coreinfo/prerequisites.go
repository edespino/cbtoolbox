package coreinfo

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// checkPrerequisites verifies that all necessary tools and configurations are available.
var checkPrerequisites = func() error {
	if err := checkGDBAvailability(); err != nil {
		return err
	}

	if _, err := getPostgresPath(); err != nil {
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
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Debug: 'file' command failed for '%s': %v\n", filePath, err)
		return false, err
	}
	return strings.Contains(string(output), "core file") || strings.Contains(string(output), "ELF"), nil
}

// validateAndAddCoreFile handles the validation of a single potential core file
// Returns true if the file is a valid core file and was added
func validateAndAddCoreFile(file string, coreFiles *[]string) error {
	valid, err := isCoreFile(file)
	if err != nil {
		return fmt.Errorf("failed to check core file %s: %v", file, err)
	}
	if valid {
		*coreFiles = append(*coreFiles, file)
	} else if verbose {
		fmt.Printf("Debug: File '%s' NOT recognized as a core file\n", file)
	}
	return nil
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
			fmt.Printf("Debug: Error accessing path '%s': %v\n", arg, err)
			continue
		}

		if info.IsDir() {
			files, err := filepath.Glob(filepath.Join(arg, "*"))
			if err != nil {
				return nil, fmt.Errorf("failed to read directory %s: %v", arg, err)
			}
			for _, file := range files {
				if err := validateAndAddCoreFile(file, &coreFiles); err != nil {
					return nil, err
				}
			}
		} else {
			if err := validateAndAddCoreFile(arg, &coreFiles); err != nil {
				return nil, err
			}
		}
	}

	if len(coreFiles) == 0 {
		return nil, fmt.Errorf("no valid core files provided")
	}
	return coreFiles, nil
}
