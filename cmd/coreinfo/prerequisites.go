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
			files, _ := filepath.Glob(filepath.Join(arg, "*"))
			for _, file := range files {
				if valid, _ := isCoreFile(file); valid {
					coreFiles = append(coreFiles, file)
				} else {
					fmt.Printf("Debug: File '%s' NOT recognized as a core file\n", file)
				}
			}
		} else {
			if valid, _ := isCoreFile(arg); valid {
				coreFiles = append(coreFiles, arg)
			} else {
				fmt.Printf("Debug: File '%s' NOT recognized as a core file\n", arg)
			}
		}
	}

	if len(coreFiles) == 0 {
		return nil, fmt.Errorf("no valid core files provided")
	}
	return coreFiles, nil
}
