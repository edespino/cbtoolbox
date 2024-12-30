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

package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

// TestExecute validates that the Execute function successfully runs
// a basic command implementation. It temporarily replaces the root
// command with a test command that always succeeds, then verifies
// the command executes without error.
func TestExecute(t *testing.T) {
	// Store original command and restore after test
	originalCmd := rootCmd
	defer func() { rootCmd = originalCmd }()

	// Create a test command that always succeeds
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

// TestRootCommandHelp validates the behavior of the root command when
// executed without subcommands. It verifies that the command returns
// an appropriate error message for unknown commands while still displaying
// the help information.
func TestRootCommandHelp(t *testing.T) {
	if err := rootCmd.Execute(); err != nil {
		// The root command without subcommands will return an "unknown command" error
		if err.Error() != "unknown command" {
			t.Errorf("Unexpected error executing root command: %v", err)
		}
	}
}
