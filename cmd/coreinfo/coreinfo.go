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

var (
	extractBasic    bool
	extractDetailed bool
	customGDBFile   string
)

// RunCoreInfo contains the logic for the coreinfo command.
func RunCoreInfo(cmd *cobra.Command, args []string) error {
	// Handle extraction
	if extractBasic {
		return extractGDBFile("gdb_commands_basic.txt", "gdb_commands_basic.txt")
	}
	if extractDetailed {
		return extractGDBFile("gdb_commands_detailed.txt", "gdb_commands_detailed.txt")
	}

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

	// Run GDB analysis
	if err := RunGDBAnalysis(coreFiles, customGDBFile); err != nil {
		return fmt.Errorf("gdb analysis failed: %v", err)
	}

	return nil
}

func init() {
	CoreinfoCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	CoreinfoCmd.Flags().BoolVarP(&extractBasic, "extract-basic", "", false, "Extract the basic GDB command file")
	CoreinfoCmd.Flags().BoolVarP(&extractDetailed, "extract-detailed", "", false, "Extract the detailed GDB command file")
	CoreinfoCmd.Flags().StringVarP(&customGDBFile, "gdb-file", "", "", "Path to a custom GDB command file")
}
