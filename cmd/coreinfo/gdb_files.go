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
