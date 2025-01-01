package coreinfo

import (
	"fmt"
	"os"
	"os/exec"
)

// RunGDBAnalysis executes the GDB analysis on the provided core files using the specified GDB command file.
func RunGDBAnalysis(coreFiles []string, gdbFile string) error {
	if gdbFile == "" {
		gdbFile = "gdb_commands_basic.txt" // Default to basic analysis
	}

	for _, coreFile := range coreFiles {
		fmt.Printf("Analyzing core file: %s using %s\n", coreFile, gdbFile)
		gdbCmd := exec.Command("gdb", "-x", gdbFile, "/usr/local/cloudberry-db/bin/postgres", coreFile)
		gdbCmd.Stdout = os.Stdout
		gdbCmd.Stderr = os.Stderr

		if err := gdbCmd.Run(); err != nil {
			return fmt.Errorf("failed to run gdb on %s: %v", coreFile, err)
		}
	}

	return nil
}
