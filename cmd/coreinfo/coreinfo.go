package coreinfo

import (
	"fmt"

	"github.com/spf13/cobra"
)

var verbose bool // Flag for verbose output

// CoreinfoCmd defines the coreinfo command for analyzing core dump files.
var CoreinfoCmd = &cobra.Command{
	Use:   "coreinfo",
	Short: "Analyze core dump files",
	Long:  "The coreinfo command analyzes core dump files to provide insights into system crashes.",
	RunE:  RunCoreInfo,
}

// RunCoreInfo contains the logic for the coreinfo command.
func RunCoreInfo(cmd *cobra.Command, args []string) error {
	// Step 1: Check prerequisites
	if err := checkPrerequisites(); err != nil {
		return fmt.Errorf("prerequisite check failed: %v", err)
	}

	// Step 2: Validate core file paths
	coreFiles, err := validateCoreFiles(args)
	if err != nil {
		return fmt.Errorf("core file validation failed: %v", err)
	}

	// Step 3: Print detailed validation results if verbose mode is enabled
	if verbose {
		for _, coreFile := range coreFiles {
			fmt.Printf("Validating file: %s -> Valid core file\n", coreFile)
		}
	}

	// Placeholder: Print core file paths (replace with actual logic later)
	fmt.Printf("Validated core files: %v\n", coreFiles)

	return nil
}

func init() {
	CoreinfoCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
}
