package coreinfo

import (
	"embed"
	"fmt"
	"os"
)

//go:embed resources/gdb_commands_basic.txt resources/gdb_commands_detailed.txt
var gdbFiles embed.FS

func extractGDBFile(filename string, outputPath string) error {
	data, err := gdbFiles.ReadFile("resources/" + filename)
	if err != nil {
		return fmt.Errorf("failed to read embedded file %s: %v", filename, err)
	}

	err = os.WriteFile(outputPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %v", outputPath, err)
	}

	fmt.Printf("File %s extracted to %s\n", filename, outputPath)
	return nil
}

func extractEmbeddedFile(filename string) (string, error) {
    // Open the embedded file
    fileContent, err := gdbFiles.ReadFile("resources/" + filename)
    if err != nil {
	return "", fmt.Errorf("failed to read embedded file %s: %v", filename, err)
    }

    // Write to a temporary file
    tmpFile, err := os.CreateTemp("", filename+"_*.txt")
    if err != nil {
	return "", fmt.Errorf("failed to create temporary file for %s: %v", filename, err)
    }

    if _, err := tmpFile.Write(fileContent); err != nil {
	return "", fmt.Errorf("failed to write content to temporary file: %v", err)
    }

    tmpFile.Close()
    return tmpFile.Name(), nil
}
