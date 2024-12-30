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

// cbtoolbox is a command-line utility for managing and monitoring
// Apache Cloudberry (Incubating) installations. It provides various
// tools for system diagnostics, configuration management, and cluster
// administration.
//
// The utility is designed to work with Apache Cloudberry (Incubating)
// installations and requires the GPHOME environment variable to be
// set to the database installation directory.
//
// Basic usage:
//
//	cbtoolbox [command] [flags]
//
// Available commands:
//   - sysinfo: Display system and database environment information
//   - help: Display help information about available commands
//
// For detailed command usage, run:
//
//	cbtoolbox help [command]
package main

import (
	"fmt"
	"os"

	"github.com/edespino/cbtoolbox/cmd"
)

// exitFunc allows for mocking the os.Exit function during testing.
// In production, it points to os.Exit.
var exitFunc = os.Exit

// run executes the root command and handles error propagation.
// It returns an error if command execution fails, after ensuring
// proper exit code is set through exitFunc.
func run() error {
	err := cmd.Execute()
	if err != nil {
		exitFunc(1)
		return err
	}
	return nil
}

// main is the entry point for the cbtoolbox utility.
// It executes the root command and handles error output.
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
