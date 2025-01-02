package coreinfo

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// checkPrerequisites verifies that all necessary tools and configurations are available.
var checkPrerequisites = func() error {
	if err := checkGDBAvailability(); err != nil {
		return fmt.Errorf("gdb not found: please install GDB using your system package manager (e.g. 'yum install gdb' or 'apt-get install gdb')")
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

// prerequisites.go
type FileInfo struct {
	Platform    string
	RealUID    string
	EffUID     string
	RealGID    string
	EffGID     string
	ExecPath   string
}

func isCoreFile(filePath string) (bool, *FileInfo, error) {
	cmd := exec.Command("file", filePath)
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Debug: 'file' command failed for '%s': %v\n", filePath, err)
		return false, nil, err
	}
	outputStr := string(output)
	isCore := strings.Contains(outputStr, "core file") || strings.Contains(outputStr, "ELF")

	var info *FileInfo
	if isCore {
		info = &FileInfo{}
		// Platform
		if match := regexp.MustCompile(`platform: '([^']+)'`).FindStringSubmatch(outputStr); len(match) > 1 {
			info.Platform = match[1]
		}
		// UIDs
		if match := regexp.MustCompile(`real uid: (\d+)`).FindStringSubmatch(outputStr); len(match) > 1 {
			info.RealUID = match[1]
		}
		if match := regexp.MustCompile(`effective uid: (\d+)`).FindStringSubmatch(outputStr); len(match) > 1 {
			info.EffUID = match[1]
		}
		// GIDs
		if match := regexp.MustCompile(`real gid: (\d+)`).FindStringSubmatch(outputStr); len(match) > 1 {
			info.RealGID = match[1]
		}
		if match := regexp.MustCompile(`effective gid: (\d+)`).FindStringSubmatch(outputStr); len(match) > 1 {
			info.EffGID = match[1]
		}
		// Executable path
		if match := regexp.MustCompile(`execfn: '([^']+)'`).FindStringSubmatch(outputStr); len(match) > 1 {
			info.ExecPath = match[1]
		}
	}

	return isCore, info, nil
}

// validateAndAddCoreFile handles the validation of a single potential core file
// Returns error if validation fails
func validateAndAddCoreFile(file string, coreFiles *[]string, coreInfos map[string]*FileInfo) error {
	valid, info, err := isCoreFile(file)
	if err != nil {
		return fmt.Errorf("failed to check core file %s: %v", file, err)
	}
	if valid {
		*coreFiles = append(*coreFiles, file)
		coreInfos[file] = info
	} else if verbose {
		fmt.Printf("Debug: File '%s' NOT recognized as a core file\n", file)
	}
	return nil
}

// validateCoreFiles validates the input paths to determine if they are core files or directories containing core files.
func validateCoreFiles(args []string) ([]string, map[string]*FileInfo, error) {
	if len(args) == 0 {
		return nil, nil, fmt.Errorf("no core files specified: usage 'cbtoolbox coreinfo <path-to-core-file>' or 'cbtoolbox coreinfo <directory-with-cores>'")
	}

	var coreFiles []string
	coreInfos := make(map[string]*FileInfo)

	for _, arg := range args {
		info, err := os.Stat(arg)
		if err != nil {
			fmt.Printf("Debug: Error accessing path '%s': %v\n", arg, err)
			continue
		}

		if info.IsDir() {
			files, err := filepath.Glob(filepath.Join(arg, "*"))
			if err != nil {
				return nil, nil, fmt.Errorf("failed to read directory %s: %v", arg, err)
			}
			for _, file := range files {
				if err := validateAndAddCoreFile(file, &coreFiles, coreInfos); err != nil {
					return nil, nil, err
				}
			}
		} else {
			if err := validateAndAddCoreFile(arg, &coreFiles, coreInfos); err != nil {
				return nil, nil, err
			}
		}
	}

	if len(coreFiles) == 0 {
		return nil, nil, fmt.Errorf("no valid core files provided")
	}
	return coreFiles, coreInfos, nil
}
