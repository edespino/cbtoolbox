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

// Package cmd implements the command-line interface for the Apache
// Cloudberry toolbox. It provides the root command and coordinates
// subcommand registration for various database management and
// monitoring functionalities.
//
// The package uses the cobra framework for command-line parsing and execution,
// organizing commands in a hierarchical structure with the root command as the
// entry point for all subcommands.
//
// Example usage:
//
//	import "github.com/edespino/cbtoolbox/cmd"
//
//	func main() {
//	    if err := cmd.Execute(); err != nil {
//	        // Handle error
//	    }
//	}
package cmd

import (
	"github.com/edespino/cbtoolbox/cmd/coreinfo"
        "github.com/edespino/cbtoolbox/cmd/sysinfo"
        "github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands.
// This command provides a help message and serves as the entry point for
// executing subcommands within the `cbtoolbox` CLI.
var rootCmd = &cobra.Command{
        Use:   "cbtoolbox",
        Short: "An Apache Cloudberry (Incubator) toolbox",
        Long:  "An Apache Cloudberry (Incubator) toolbox",
}

// init registers all subcommands with the root command.
// Currently registered commands:
//   - sysinfo: Displays system and database environment information
//   - coreinfo: Analyzes core dump files for diagnostic purposes
func init() {
        rootCmd.AddCommand(sysinfo.Cmd)
	rootCmd.AddCommand(coreinfo.CoreinfoCmd)

}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
// Returns an error if command execution fails.
func Execute() error {
        return rootCmd.Execute()
}
