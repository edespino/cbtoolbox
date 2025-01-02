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
// root.go

package cmd

import (
        "fmt"
        "os"

        "github.com/edespino/cbtoolbox/cmd/coreinfo"
        "github.com/edespino/cbtoolbox/cmd/sysinfo"
        "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
        Use:   "cbtoolbox",
        Short: "An Apache Cloudberry (Incubator) toolbox",
        Long:  "An Apache Cloudberry (Incubator) toolbox",
        PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
                // Skip GPHOME check for help and version commands
                if cmd.Name() == "help" || cmd.Name() == "version" {
                        return nil
                }

                // Skip check if this is the root command being executed without subcommands
                if cmd.Name() == "cbtoolbox" {
                        return nil
                }

                // Check GPHOME environment variable
                gphome := os.Getenv("GPHOME")
                if gphome == "" {
                        return fmt.Errorf("GPHOME environment variable is not set")
                }

                // Verify GPHOME points to a valid directory
                if _, err := os.Stat(gphome); os.IsNotExist(err) {
                        return fmt.Errorf("GPHOME directory does not exist: %s", gphome)
                }

                return nil
        },
}

func init() {
        rootCmd.AddCommand(sysinfo.Cmd)
        rootCmd.AddCommand(coreinfo.CoreinfoCmd)
}

func Execute() error {
        return rootCmd.Execute()
}
