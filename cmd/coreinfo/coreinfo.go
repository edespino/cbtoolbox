package coreinfo

import (
	"fmt"

	"github.com/spf13/cobra"
)

// CoreinfoCmd represents the coreinfo command
var CoreinfoCmd = &cobra.Command{
	Use:   "coreinfo",
	Short: "Analyze core dump files",
	Long:  "The coreinfo command analyzes core dump files to provide insights into system crashes.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Step 1: Check prerequisites
		if err := checkPrerequisites(); err != nil {
			return fmt.Errorf("prerequisite check failed: %v", err)
		}

		// Step 2: Validate core file paths
		coreFiles, err := validateCoreFiles(args)
		if err != nil {
			return fmt.Errorf("core file validation failed: %v", err)
		}

		// Placeholder: Print core file paths (replace with actual logic later)
		fmt.Printf("Validated core files: %v\n", coreFiles)

		return nil
	},
}
