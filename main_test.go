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

package main

import (
	"os"
	"strings"
	"testing"
)

// TestRun validates the main command execution path of cbtoolbox.
// It tests various command scenarios including help display,
// invalid command handling, and required environment variable validation.
func TestRun(t *testing.T) {
	// Store original state to restore after tests
	originalExit := exitFunc
	originalArgs := os.Args
	originalGPHOME := os.Getenv("GPHOME")

	defer func() {
		exitFunc = originalExit
		os.Args = originalArgs
		os.Setenv("GPHOME", originalGPHOME)
	}()

	// Test cases cover core functionality and error conditions
	tests := []struct {
		name     string   // Test case description
		args     []string // Command line arguments
		gphome   string   // GPHOME environment value
		wantErr  bool     // Whether an error is expected
		errMsg   string   // Expected error message substring
		wantExit bool     // Whether os.Exit should be called
	}{
		{
			name:     "help command",
			args:     []string{"cbtoolbox", "--help"},
			wantErr:  false,
			wantExit: false,
		},
		{
			name:     "invalid command",
			args:     []string{"cbtoolbox", "invalid"},
			wantErr:  true,
			errMsg:   "unknown command",
			wantExit: true,
		},
		{
			name:     "sysinfo without GPHOME",
			args:     []string{"cbtoolbox", "sysinfo"},
			gphome:   "",
			wantErr:  true,
			errMsg:   "GPHOME environment",
			wantExit: true,
		},
	}

	// Execute each test case in isolation
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exitCalled := false
			exitFunc = func(code int) {
				exitCalled = true
				if code != 1 {
					t.Errorf("Expected exit code 1, got %d", code)
				}
			}

			os.Args = tt.args
			os.Setenv("GPHOME", tt.gphome)

			err := run()

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Expected error containing %q, got %v", tt.errMsg, err)
				}
				if tt.wantExit && !exitCalled {
					t.Error("Expected exit function to be called")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if exitCalled {
					t.Error("Exit function called unexpectedly")
				}
			}
		})
	}
}
