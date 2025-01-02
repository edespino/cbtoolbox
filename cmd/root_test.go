// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// root_test.go
package cmd

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
)

func TestExecute(t *testing.T) {
	originalCmd := rootCmd
	defer func() { rootCmd = originalCmd }()

	rootCmd = &cobra.Command{
		Use:   "test",
		Short: "Test command",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	if err := Execute(); err != nil {
		t.Errorf("Execute() failed: %v", err)
	}
}

func TestRootCommandHelp(t *testing.T) {
	if err := rootCmd.Execute(); err != nil {
		if err.Error() != "unknown command" {
			t.Errorf("Unexpected error executing root command: %v", err)
		}
	}
}

func TestGPHOMEValidation(t *testing.T) {
	originalGPHOME := os.Getenv("GPHOME")
	defer os.Setenv("GPHOME", originalGPHOME)

	tests := []struct {
		name        string
		gphomePath  string
		createDir   bool
		shouldError bool
	}{
		{
			name:        "GPHOME not set",
			gphomePath:  "",
			createDir:   false,
			shouldError: true,
		},
		{
			name:        "GPHOME set to non-existent directory",
			gphomePath:  "/nonexistent/path",
			createDir:   false,
			shouldError: true,
		},
		{
			name:        "GPHOME set to valid directory",
			gphomePath:  t.TempDir(),
			createDir:   true,
			shouldError: false,
		},
	}

	testCmd := &cobra.Command{
		Use: "test",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	rootCmd.AddCommand(testCmd)
	defer rootCmd.RemoveCommand(testCmd)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.gphomePath == "" {
				os.Unsetenv("GPHOME")
			} else {
				os.Setenv("GPHOME", tt.gphomePath)
			}

			err := rootCmd.PersistentPreRunE(testCmd, []string{})

			if (err != nil) != tt.shouldError {
				t.Errorf("PersistentPreRunE() error = %v, shouldError = %v", err, tt.shouldError)
			}
		})
	}
}

func TestGPHOMESkipForHelp(t *testing.T) {
	originalGPHOME := os.Getenv("GPHOME")
	os.Unsetenv("GPHOME")
	defer os.Setenv("GPHOME", originalGPHOME)

	helpCmd := &cobra.Command{
		Use: "help",
	}

	if err := rootCmd.PersistentPreRunE(helpCmd, []string{}); err != nil {
		t.Errorf("PersistentPreRunE() should not check GPHOME for help command, got error: %v", err)
	}
}
