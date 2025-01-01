package coreinfo

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

// RunGDBAnalysisWithSummary performs GDB analysis and includes a summary at the top of the output.
func RunGDBAnalysisWithSummary(coreFiles []string, customGDBFile string) error {
	for _, coreFile := range coreFiles {
		// Select GDB file
		gdbFile := "gdb_commands_basic.txt" // Default to basic commands
		if customGDBFile != "" {
			gdbFile = customGDBFile
		}

		// Run GDB command
		gdbCmd := exec.Command("gdb", "-q", "-x", gdbFile, "/usr/local/cloudberry-db/bin/postgres", coreFile)
		output, err := gdbCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to run GDB on %s: %v", coreFile, err)
		}

		// Extract and print summary
		summary, err := extractCoreSummary(string(output))
		if err != nil {
			return fmt.Errorf("failed to extract core summary for %s: %v", coreFile, err)
		}
		fmt.Println(summary)

		// Print the full GDB output after the summary
		fmt.Println("Detailed GDB Output:")
		fmt.Println(string(output))
	}

	return nil
}

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

func extractCoreSummary(gdbOutput string) (string, error) {
    var binary, signal, faultAddr, threadID, processArgs string

    // Correctly escaped regex patterns
    binaryRegex := regexp.MustCompile(`Core was generated by '([^']+)'`)
    signalRegex := regexp.MustCompile(`Program terminated with signal (\w+), (.+)`)
    faultAddrRegex := regexp.MustCompile(`si_addr = ([^,]+)`)
    threadIDRegex := regexp.MustCompile(`Current thread is (\d+)`)
    argsRegex := regexp.MustCompile(`Core was generated by '([^']+)'.*`)

    // Match and extract relevant information
    if match := binaryRegex.FindStringSubmatch(gdbOutput); len(match) > 1 {
	binary = match[1]
    } else {
	return "", fmt.Errorf("failed to extract binary information")
    }

    if match := signalRegex.FindStringSubmatch(gdbOutput); len(match) > 2 {
	signal = fmt.Sprintf("%s (%s)", match[1], match[2])
    } else {
	signal = "Unknown signal"
    }

    if match := faultAddrRegex.FindStringSubmatch(gdbOutput); len(match) > 1 {
	faultAddr = match[1]
    } else {
	faultAddr = "N/A"
    }

    if match := threadIDRegex.FindStringSubmatch(gdbOutput); len(match) > 1 {
	threadID = match[1]
    } else {
	threadID = "N/A"
    }

    if match := argsRegex.FindStringSubmatch(gdbOutput); len(match) > 1 {
	processArgs = match[1]
    } else {
	processArgs = "N/A"
    }

    // Format the summary
    summary := fmt.Sprintf(`
Core Dump Analysis Summary:
----------------------------------------
- Binary: %s
- Signal: %s
- Faulting Address: %s
- Thread ID: %s
- Process Args: %s
`, binary, signal, faultAddr, threadID, processArgs)

    return summary, nil
}
