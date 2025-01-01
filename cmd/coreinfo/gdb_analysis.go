package coreinfo

import (
	"fmt"
	"os"
	"os/exec"
)

// RunGDBAnalysis uses the embedded GDB command files for analysis.
func RunGDBAnalysis(coreFiles []string, gdbFile string) error {
	if gdbFile == "" {
		// Use embedded basic GDB commands by default
		gdbFile = "gdb_commands_basic.txt"
		tempFile, err := os.CreateTemp("", gdbFile)
		if err != nil {
			return fmt.Errorf("failed to create temp file for GDB commands: %v", err)
		}
		defer tempFile.Close()
		defer os.Remove(tempFile.Name()) // Clean up after use

		// Write the embedded file to the temporary location
		data, err := gdbFiles.ReadFile("resources/" + gdbFile)
		if err != nil {
			return fmt.Errorf("failed to read embedded GDB commands: %v", err)
		}
		if _, err := tempFile.Write(data); err != nil {
			return fmt.Errorf("failed to write GDB commands to temp file: %v", err)
		}
		gdbFile = tempFile.Name() // Update to the temp file path
	}

	for _, coreFile := range coreFiles {
		fmt.Printf("Analyzing core file: %s using %s\n", coreFile, gdbFile)

		// Construct the GDB command with the --quiet option
		gdbCmd := exec.Command("gdb", "--quiet", "-x", gdbFile, "/usr/local/cloudberry-db/bin/postgres", coreFile)

		// Redirect GDB output to the standard output and error streams
		gdbCmd.Stdout = os.Stdout
		gdbCmd.Stderr = os.Stderr

		// Execute the GDB command and handle errors
		if err := gdbCmd.Run(); err != nil {
			return fmt.Errorf("failed to run GDB on %s: %v", coreFile, err)
		}
	}

	return nil
}
