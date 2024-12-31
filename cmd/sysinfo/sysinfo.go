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

// Package sysinfo implements functionality to gather and display detailed
// system and database environment information for Apache Cloudberry installations.
//
// This package requires the GPHOME environment variable to be set for database-specific
// information. If GPHOME is not set, it provides only system-level details.
//
// Supported output formats:
//   - YAML (default)
//   - JSON
//
// Example usage:
//
//	import "github.com/edespino/cbtoolbox/cmd/sysinfo"
//
//	cmd := sysinfo.Cmd
//	err := cmd.Execute()
package sysinfo

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// Package-level variables that control behavior and configuration.
var (
	// formatFlag determines the output format (yaml or json)
	formatFlag string

	// procMeminfo specifies the path to system memory information
	procMeminfo = "/proc/meminfo"
	osReleasePath = "/etc/os-release"
)

// Cmd represents the sysinfo command that gathers and displays
// system and database environment information.
var Cmd = &cobra.Command{
	Use:   "sysinfo",
	Short: "Display system information",
	Long: `Gather and display detailed system and database environment information.
Requires GPHOME environment variable to be set for database-specific information.`,
	RunE: RunSysInfo,
}

// SysInfo represents the complete system and database environment
// information collected by the sysinfo command.
type SysInfo struct {
	OS                 string            `json:"os" yaml:"os"`
	Architecture       string            `json:"architecture" yaml:"architecture"`
	Hostname           string            `json:"hostname" yaml:"hostname"`
	Kernel             string            `json:"kernel" yaml:"kernel"`
	OSVersion          string            `json:"os_version" yaml:"os_version"`
	CPUs               int               `json:"cpus" yaml:"cpus"`
	MemoryStats        map[string]string `json:"memory_stats" yaml:"memory_stats"`
	GPHOME             string            `json:"GPHOME,omitempty" yaml:"GPHOME,omitempty"`
	PGConfigConfigure  []string          `json:"pg_config_configure,omitempty" yaml:"pg_config_configure,omitempty"`
	PostgresVersion    string            `json:"postgres_version,omitempty" yaml:"postgres_version,omitempty"`
	GPVersion          string            `json:"gp_version,omitempty" yaml:"gp_version,omitempty"`
}

// init initializes the sysinfo command configuration.
// It sets up the default output format and command flags.
func init() {
	// Default output format is YAML
	formatFlag = "yaml"
	Cmd.Flags().StringVar(&formatFlag, "format", "yaml", "Output format: yaml or json")
}

// validateFormat checks if the provided format is supported.
// Returns nil for valid formats (yaml, json) and an error for unsupported formats.
func validateFormat(format string) error {
	switch format {
	case "yaml", "json":
		return nil
	default:
		return fmt.Errorf("invalid format: %s (supported formats: yaml, json)", format)
	}
}

// readFile abstracts file reading logic, making it mockable during tests.
var readFile = os.ReadFile

// getOS returns the operating system name using runtime information.
// This function provides a consistent way to determine the OS across different platforms.
func getOS() string {
	return runtime.GOOS
}

// getArchitecture returns the system's CPU architecture using runtime information.
// This provides the underlying hardware architecture of the system.
func getArchitecture() string {
	return runtime.GOARCH
}

// getHostname returns the system's network hostname.
// Returns an error if the hostname cannot be determined.
func getHostname() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("hostname: failed to retrieve hostname: %w", err)
	}
	return hostname, nil
}

// getKernelVersion returns the Linux kernel version by executing 'uname -r'.
// Returns an error if the command fails or cannot be executed.
func getKernelVersion() (string, error) {
	output, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return "", fmt.Errorf("kernel: failed to retrieve version: %w", err)
	}
	return "Linux " + strings.TrimSpace(string(output)), nil
}

// getOSVersion returns the operating system version from /etc/os-release.
// Extracts the PRETTY_NAME field from the file.
// Returns "unknown" if the PRETTY_NAME field is not found.
func getOSVersion() (string, error) {
	content, err := readFile(osReleasePath)
	if err != nil {
		return "", fmt.Errorf("os-release: failed to read file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			return strings.Trim(line[len("PRETTY_NAME="):], `"`), nil
		}
	}
	return "", fmt.Errorf("os-release: PRETTY_NAME not found")
}

// getCPUCount returns the number of CPU cores available to the system
// using runtime information.
func getCPUCount() int {
	return runtime.NumCPU()
}

// getReadableMemoryStats returns memory statistics from /proc/meminfo in a human-readable format.
// The returned map includes MemTotal, MemFree, MemAvailable, Cached, and Buffers,
// with values converted to appropriate units (KiB, MiB, GiB).
func getReadableMemoryStats() (map[string]string, error) {
	output, err := os.ReadFile(procMeminfo)
	if err != nil {
		return nil, fmt.Errorf("meminfo: failed to read file: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	memoryStats := make(map[string]string)
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		key := strings.TrimSuffix(fields[0], ":")
		value := fields[1]
		if key == "MemTotal" || key == "MemFree" || key == "MemAvailable" || key == "Cached" || key == "Buffers" {
			converted := humanizeSize(value)
			memoryStats[key] = converted
		}
	}
	return memoryStats, nil
}

// humanizeSize converts a memory size from kilobytes to a human-readable string.
// Input is a string representing kilobytes.
// Output format:
//   - For values >= 1024*1024 KB: X.X GiB
//   - For values >= 1024 KB: X.X MiB
//   - For values < 1024 KB: X KiB
func humanizeSize(kb string) string {
	kbInt, err := strconv.Atoi(kb)
	if err != nil {
		return kb
	}
	switch {
	case kbInt >= 1024*1024:
		return fmt.Sprintf("%.1f GiB", float64(kbInt)/(1024*1024))
	case kbInt >= 1024:
		return fmt.Sprintf("%.1f MiB", float64(kbInt)/1024)
	default:
		return fmt.Sprintf("%d KiB", kbInt)
	}
}

// getGPHOME returns and validates the GPHOME environment variable.
// Returns the GPHOME path if it exists and is valid.
// Returns an error if:
//   - GPHOME environment variable is not set
//   - GPHOME directory does not exist
func getGPHOME() (string, error) {
	gphome := os.Getenv("GPHOME")
	if gphome == "" {
		return "", fmt.Errorf("GPHOME: environment variable not set")
	}
	if _, err := os.Stat(gphome); os.IsNotExist(err) {
		return gphome, fmt.Errorf("GPHOME: directory does not exist: %s", gphome)
	}
	return gphome, nil
}

// getPGConfigConfigure returns PostgreSQL build configuration options.
// Executes pg_config --configure in the specified GPHOME/bin directory.
// Returns an error if:
//   - pg_config executable is not found in GPHOME/bin
//   - pg_config command execution fails
func getPGConfigConfigure(gphome string) ([]string, error) {
	pgConfigPath := filepath.Join(gphome, "bin", "pg_config")
	if _, err := os.Stat(pgConfigPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("pg_config: file not found at %s", pgConfigPath)
	}

	cmd := exec.Command(pgConfigPath, "--configure")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("pg_config: failed to execute: %w", err)
	}
	config := strings.ReplaceAll(strings.TrimSpace(string(output)), "'", "")
	return strings.Fields(config), nil
}

// getPostgresVersion returns the PostgreSQL server version.
// Executes postgres --version in the specified GPHOME/bin directory.
// Returns an error if:
//   - postgres executable is not found in GPHOME/bin
//   - postgres command execution fails
func getPostgresVersion(gphome string) (string, error) {
	postgresPath := filepath.Join(gphome, "bin", "postgres")
	if _, err := os.Stat(postgresPath); os.IsNotExist(err) {
		return "", fmt.Errorf("postgres: executable not found at %s", postgresPath)
	}

	cmd := exec.Command(postgresPath, "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("postgres: failed to execute version check: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// getGPVersion returns the Apache Cloudberry version.
// Executes postgres --gp-version in the specified GPHOME/bin directory.
// Returns an error if:
//   - postgres executable is not found in GPHOME/bin
//   - postgres command execution fails
func getGPVersion(gphome string) (string, error) {
	postgresPath := filepath.Join(gphome, "bin", "postgres")
	if _, err := os.Stat(postgresPath); os.IsNotExist(err) {
		return "", fmt.Errorf("postgres: executable not found at %s", postgresPath)
	}

	cmd := exec.Command(postgresPath, "--gp-version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("postgres: failed to execute gp-version check: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// gatherGPHOMEInfo collects all database-related information.
// Returns:
//   - string: GPHOME path if valid
//   - []string: PostgreSQL build configuration options
//   - string: PostgreSQL server version
//   - string: Apache Cloudberry version
//   - []error: Collection of any errors encountered
//
// If GPHOME is not set or invalid, returns appropriate error messages for each
// component that could not be checked.
func gatherGPHOMEInfo() (string, []string, string, string, []error) {
	gphome, gphomeErr := getGPHOME()
	var pgConfig []string
	var postgresVersion string
	var gpVersion string
	var errs []error

	if gphomeErr != nil {
		errs = append(errs, fmt.Errorf("GPHOME error: %w", gphomeErr))
	}

	if gphome != "" {
		config, err := getPGConfigConfigure(gphome)
		if err != nil {
			errs = append(errs, fmt.Errorf("pg_config error: %w", err))
		} else {
			pgConfig = config
		}

		version, err := getPostgresVersion(gphome)
		if err != nil {
			errs = append(errs, fmt.Errorf("postgres version error: %w", err))
		} else {
			postgresVersion = version
		}

		gpVer, err := getGPVersion(gphome)
		if err != nil {
			errs = append(errs, fmt.Errorf("gp version error: %w", err))
		} else {
			gpVersion = gpVer
		}
	}

	return gphome, pgConfig, postgresVersion, gpVersion, errs
}

// RunSysInfo gathers and displays system and database information.
// This is the main entry point for the sysinfo command execution.
//
// The function performs the following steps:
//  1. Validates the output format (yaml/json)
//  2. Checks for GPHOME environment variable
//  3. Gathers system information concurrently
//  4. Collects database information if GPHOME is set
//  5. Formats and displays the collected information
//
// Returns an error if:
//   - The format is invalid
//   - Required system information cannot be collected
//   - GPHOME is not set (after displaying available system information)
func RunSysInfo(cmd *cobra.Command, args []string) error {
	if err := validateFormat(formatFlag); err != nil {
		return err
	}

	// Check GPHOME first
	if os.Getenv("GPHOME") == "" {
		info := SysInfo{
			OS:           getOS(),
			Architecture: getArchitecture(),
			CPUs:         getCPUCount(),
		}

		// Get other system info
		if hostname, err := getHostname(); err == nil {
			info.Hostname = hostname
		}
		if kernel, err := getKernelVersion(); err == nil {
			info.Kernel = kernel
		}
		if osVersion, err := getOSVersion(); err == nil {
			info.OSVersion = osVersion
		}
		if memStats, err := getReadableMemoryStats(); err == nil {
			info.MemoryStats = memStats
		}

		// Output the available information
		var output []byte
		var err error
		if formatFlag == "json" {
			output, err = json.MarshalIndent(info, "", "  ")
		} else {
			output, err = yaml.Marshal(info)
		}
		if err != nil {
			return fmt.Errorf("output: failed to generate: %w", err)
		}

		fmt.Println(string(output))
		return fmt.Errorf("GPHOME environment variable is not set")
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	info := SysInfo{}
	errs := make([]error, 0)

	// Concurrent data collection for system information
	wg.Add(7)
	go func() { defer wg.Done(); info.OS = getOS() }()
	go func() { defer wg.Done(); info.Architecture = getArchitecture() }()
	go func() {
		defer wg.Done()
		if hostname, err := getHostname(); err == nil {
			info.Hostname = hostname
		} else {
			mu.Lock()
			errs = append(errs, err)
			mu.Unlock()
		}
	}()
	go func() {
		defer wg.Done()
		if kernel, err := getKernelVersion(); err == nil {
			info.Kernel = kernel
		} else {
			mu.Lock()
			errs = append(errs, err)
			mu.Unlock()
		}
	}()
	go func() {
		defer wg.Done()
		if osVersion, err := getOSVersion(); err == nil {
			info.OSVersion = osVersion
		} else {
			mu.Lock()
			errs = append(errs, err)
			mu.Unlock()
		}
	}()
	go func() { defer wg.Done(); info.CPUs = getCPUCount() }()
	go func() {
		defer wg.Done()
		if memStats, err := getReadableMemoryStats(); err == nil {
			mu.Lock()
			info.MemoryStats = memStats
			mu.Unlock()
		} else {
			mu.Lock()
			info.MemoryStats = map[string]string{"error": err.Error()}
			errs = append(errs, err)
			mu.Unlock()
		}
	}()

	// Collect database-specific information
	gphome, pgConfig, postgresVersion, gpVersion, gphomeErrs := gatherGPHOMEInfo()
	if gphome != "" {
		info.GPHOME = gphome
		info.PGConfigConfigure = pgConfig
		info.PostgresVersion = postgresVersion
		info.GPVersion = gpVersion
	}

	wg.Wait()

	// Handle and report any errors that occurred during collection
	if len(errs) > 0 || len(gphomeErrs) > 0 {
		fmt.Println("\nSummary of errors:")
		for _, err := range errs {
			fmt.Println("-", err)
		}
		for _, err := range gphomeErrs {
			fmt.Println("-", err)
		}

		// Only fail if we have errors from required components
		if len(errs) > 0 || len(gphomeErrs) > 0 {
			return fmt.Errorf("errors occurred during system info collection")
		}
	}

	// Generate output in requested format
	var output []byte
	var err error
	if formatFlag == "json" {
		output, err = json.MarshalIndent(info, "", "  ")
	} else {
		output, err = yaml.Marshal(info)
	}
	if err != nil {
		return fmt.Errorf("output: failed to generate: %w", err)
	}

	fmt.Println(string(output))
	return nil
}
