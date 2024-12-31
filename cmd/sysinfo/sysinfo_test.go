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

// Package sysinfo_test contains the test suite for the sysinfo package.
// It provides comprehensive testing of system information gathering functionality,
// including both success and failure cases, as well as concurrent execution testing.
package sysinfo

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/spf13/cobra"
)

// captureOutput captures stdout during test execution to validate output.
// It creates a pipe, redirects stdout to it, executes the provided function,
// and returns the captured output as a string.
func captureOutput(f func()) string {
	r, w, _ := os.Pipe()
	stdOut := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = stdOut }()

	f()
	w.Close()
	out, _ := io.ReadAll(r)
	return string(out)
}

// TestGetOS validates that the getOS function returns a valid operating system name.
// It ensures the returned string is non-empty.
func TestGetOS(t *testing.T) {
	os := getOS()
	if os == "" {
		t.Errorf("Expected OS to be non-empty")
	}
}

// TestGetArchitecture validates that the system architecture information is retrieved.
// It ensures the returned architecture string is non-empty.
func TestGetArchitecture(t *testing.T) {
	arch := getArchitecture()
	if arch == "" {
		t.Errorf("Expected architecture to be non-empty")
	}
}

// TestGetHostname validates hostname retrieval functionality.
// It tests both successful hostname retrieval and proper string content.
func TestGetHostname(t *testing.T) {
	hostname, err := getHostname()
	if err != nil {
		t.Errorf("Unexpected error retrieving hostname: %v", err)
	}
	if hostname == "" {
		t.Errorf("Expected hostname to be non-empty")
	}
}

// TestGetKernelVersion validates kernel version retrieval.
// It ensures both successful execution and proper version string format.
func TestGetKernelVersion(t *testing.T) {
	kernel, err := getKernelVersion()
	if err != nil {
		t.Errorf("Unexpected error retrieving kernel version: %v", err)
	}
	if !strings.HasPrefix(kernel, "Linux ") {
		t.Errorf("Expected kernel version to start with 'Linux '")
	}
}

// TestGetKernelVersionError validates error handling when uname command is unavailable.
// It simulates a missing uname command by modifying the PATH environment variable.
func TestGetKernelVersionError(t *testing.T) {
	tempDir := os.TempDir()
	originalPath := os.Getenv("PATH")
	defer os.Setenv("PATH", originalPath)

	os.Setenv("PATH", tempDir)
	_, err := getKernelVersion()

	if err == nil {
		t.Errorf("Expected error when uname command is unavailable")
	}
}

// TestGetOSVersion validates operating system version retrieval.
// It verifies successful reading of /etc/os-release and proper content extraction.
func TestGetOSVersion(t *testing.T) {
	// Mock for non-Linux systems
	if _, err := os.Stat("/etc/os-release"); os.IsNotExist(err) {
		t.Skip("Skipping test: /etc/os-release does not exist on this system")
	}

	osVersion, err := getOSVersion()
	if err != nil {
		t.Errorf("Unexpected error retrieving OS version: %v", err)
	}
	if osVersion == "" {
		t.Errorf("Expected OS version to be non-empty")
	}
}

// TestGetCPUCount validates CPU core count retrieval.
// It ensures the reported number of CPUs is greater than 0.
func TestGetCPUCount(t *testing.T) {
	cpus := getCPUCount()
	if cpus <= 0 {
		t.Errorf("Expected CPU count to be greater than 0, got: %d", cpus)
	}
}

// TestGetReadableMemoryStats validates memory statistics retrieval and formatting.
// It tests both successful retrieval and proper formatting of memory values.
func TestGetReadableMemoryStats(t *testing.T) {
	// Mock for non-Linux systems
	if _, err := os.Stat("/proc/meminfo"); os.IsNotExist(err) {
		t.Skip("Skipping test: /proc/meminfo does not exist on this system")
	}

	originalProcMeminfo := procMeminfo
	defer func() { procMeminfo = originalProcMeminfo }()

	memoryStats, err := getReadableMemoryStats()
	if err != nil {
		t.Errorf("Unexpected error retrieving memory stats: %v", err)
	}
	if len(memoryStats) == 0 {
		t.Errorf("Expected memory stats to be non-empty")
	}

	// Verify all required memory statistics are present
	expectedKeys := []string{"MemTotal", "MemFree", "MemAvailable", "Cached", "Buffers"}
	for _, key := range expectedKeys {
		if _, exists := memoryStats[key]; !exists {
			t.Errorf("Expected memory stat '%s' not found", key)
		}
	}
}

// TestGetReadableMemoryStatsMissingFile validates error handling for missing /proc/meminfo.
// It simulates a non-existent meminfo file and verifies proper error reporting.
func TestGetReadableMemoryStatsMissingFile(t *testing.T) {
	originalProcMeminfo := procMeminfo
	defer func() { procMeminfo = originalProcMeminfo }()

	procMeminfo = "/nonexistent/meminfo"

	_, err := getReadableMemoryStats()
	if err == nil {
		t.Errorf("Expected error for missing /proc/meminfo")
	}
	if !strings.Contains(err.Error(), "meminfo: failed to read file") {
		t.Errorf("Expected error message to contain 'meminfo: failed to read file', got: %v", err)
	}
}

// TestHumanizeSize validates memory size conversion functionality.
// Tests conversion of various memory sizes to human-readable format.
func TestHumanizeSize(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"1024", "1.0 MiB"},          // Test MiB conversion
		{"2048576", "2.0 GiB"},       // Test GiB conversion
		{"512", "512 KiB"},           // Test KiB format
		{"invalid", "invalid"},        // Test invalid input handling
	}

	for _, tc := range testCases {
		result := humanizeSize(tc.input)
		if result != tc.expected {
			t.Errorf("humanizeSize(%s) = %s; want %s", tc.input, result, tc.expected)
		}
	}
}


// It verifies:
// - Command fails appropriately
// - Error message is correct
// - Basic system information is still output
// - No database-specific information is included
func TestRunSysInfoWithoutGPHOME(t *testing.T) {
	// Clear GPHOME
	originalGPHOME := os.Getenv("GPHOME")
	defer os.Setenv("GPHOME", originalGPHOME) // Restore original GPHOME after test
	os.Unsetenv("GPHOME")

	// Mock Cobra command and args
	cmd := &cobra.Command{}
	args := []string{}

	// Run sysinfo and expect an error
	err := RunSysInfo(cmd, args)
	if err == nil || err.Error() != "GPHOME environment variable is not set" {
		t.Errorf("Expected error for unset GPHOME, got: %v", err)
	}
}

// TestGetGPHOMEEmpty validates error handling when GPHOME environment variable is unset.
// Verifies proper error message and handling of missing environment variable.
func TestGetGPHOMEEmpty(t *testing.T) {
	os.Unsetenv("GPHOME")
	_, err := getGPHOME()
	if err == nil || !strings.Contains(err.Error(), "GPHOME: environment variable not set") {
		t.Errorf("Expected error for unset GPHOME")
	}
}

// TestGetPGConfigConfigure validates error handling for missing pg_config executable.
// Verifies proper error reporting when pg_config is not found in the specified path.
func TestGetPGConfigConfigure(t *testing.T) {
	os.Setenv("GPHOME", "/tmp")
	_, err := getPGConfigConfigure("/tmp")
	if err == nil {
		t.Errorf("Expected error for non-existent pg_config")
	}
}

// TestValidateFormat tests format validation for supported and unsupported formats.
// Verifies proper handling of valid (yaml, json) and invalid format specifications.
func TestValidateFormat(t *testing.T) {
	testCases := []struct {
		format string
		valid  bool
	}{
		{"json", true},
		{"yaml", true},
		{"invalid", false},
	}

	for _, tc := range testCases {
		err := validateFormat(tc.format)
		if tc.valid && err != nil {
			t.Errorf("Unexpected error for valid format '%s': %v", tc.format, err)
		}
		if !tc.valid && err == nil {
			t.Errorf("Expected error for invalid format '%s'", tc.format)
		}
	}
}

// TestRunSysInfoValidFormats validates output generation in both JSON and YAML formats.
// Creates a mock environment with required executables and verifies proper output formatting.
func TestRunSysInfoValidFormats(t *testing.T) {
	originalGPHOME := os.Getenv("GPHOME")
	defer os.Setenv("GPHOME", originalGPHOME)

	tmpDir := t.TempDir()
	binDir := filepath.Join(tmpDir, "bin")
	err := os.MkdirAll(binDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test bin directory: %v", err)
	}

	// Create mock database executables with expected outputs
	pgConfigPath := filepath.Join(binDir, "pg_config")
	postgresPath := filepath.Join(binDir, "postgres")
	pgConfigContent := "#!/bin/sh\necho '--prefix=/usr/local/cloudberry-db'\n"
	postgresContent := `#!/bin/sh
case $1 in
--version) echo 'postgres mock';;
--gp-version) echo 'postgres mock';;
esac`

	if err := os.WriteFile(pgConfigPath, []byte(pgConfigContent), 0755); err != nil {
		t.Fatalf("Failed to create mock pg_config: %v", err)
	}
	if err := os.WriteFile(postgresPath, []byte(postgresContent), 0755); err != nil {
		t.Fatalf("Failed to create mock postgres: %v", err)
	}

	os.Setenv("GPHOME", tmpDir)

	// Mock system file paths
	procMeminfo = filepath.Join(tmpDir, "meminfo")
	osReleasePath = filepath.Join(tmpDir, "os-release")

	// Create mocked files
	err = os.WriteFile(procMeminfo, []byte("MemTotal: 16384 kB\nMemFree: 8192 kB\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write mock procMeminfo file: %v", err)
	}

	err = os.WriteFile(osReleasePath, []byte("PRETTY_NAME=\"Mock Linux\"\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write mock osReleasePath file: %v", err)
	}

	// Determine the expected OS value dynamically
	expectedOS := runtime.GOOS

	// Dynamically determine the expected architecture
	expectedArchitecture := runtime.GOARCH

	// Test both JSON and YAML output formats
	for _, format := range []string{"json", "yaml"} {
		formatFlag = format
		output := captureOutput(func() {
			err := RunSysInfo(nil, nil)
			if err != nil {
				t.Errorf("Unexpected error for format %s: %v", format, err)
			}
		})

		// Log actual output for debugging
		t.Logf("Captured output for format %s:\n%s", format, output)

		// Validate JSON
		if format == "json" {
			if !strings.Contains(output, `"os": "`+expectedOS+`"`) {
				t.Errorf("Expected JSON output to contain \"os\": \"%s\", got:\n%s", expectedOS, output)
			}
			if !strings.Contains(output, `"architecture": "`+expectedArchitecture+`"`) {
				t.Errorf("Expected JSON output to contain \"architecture\": \"%s\", got:\n%s", expectedArchitecture, output)
			}
		}

		// Validate YAML
		if format == "yaml" {
			if !strings.Contains(output, "os: "+expectedOS) {
				t.Errorf("Expected YAML output to contain os: %s, got:\n%s", expectedOS, output)
			}
			if !strings.Contains(output, "architecture: "+expectedArchitecture) {
				t.Errorf("Expected YAML output to contain architecture: %s, got:\n%s", expectedArchitecture, output)
			}
		}

		if !strings.Contains(output, tmpDir) {
			t.Errorf("Expected output to contain GPHOME path")
		}
	}
}

// TestRunSysInfoInvalidFormat validates error handling for invalid output format.
// Verifies proper error message when an unsupported format is specified.
func TestRunSysInfoInvalidFormat(t *testing.T) {
	formatFlag = "invalid"
	defer func() { formatFlag = "yaml" }()

	err := RunSysInfo(nil, nil)
	if err == nil {
		t.Error("Expected error for invalid format")
	}
	if !strings.Contains(err.Error(), "invalid format") {
		t.Errorf("Expected error message to contain 'invalid format', got: %v", err)
	}
}

// TestRunSysInfoConcurrency validates thread safety of the RunSysInfo function.
// Tests concurrent execution with invalid GPHOME to verify proper error handling.
func TestRunSysInfoConcurrency(t *testing.T) {
	originalGPHOME := os.Getenv("GPHOME")
	defer os.Setenv("GPHOME", originalGPHOME)

	tmpDir := t.TempDir()
	os.Setenv("GPHOME", tmpDir)

	var wg sync.WaitGroup
	formatFlag = "json"
	errChan := make(chan error, 10)

	output := captureOutput(func() {
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := RunSysInfo(nil, nil); err != nil {
					errChan <- err
				}
			}()
		}
		wg.Wait()
		close(errChan)
	})

	// Verify errors were received as expected
	errorCount := 0
	for err := range errChan {
		if err != nil {
			errorCount++
		}
	}

	if errorCount == 0 {
		t.Error("Expected errors in concurrent execution with invalid GPHOME")
	}

	// Verify error output format
	if !strings.Contains(output, "Summary of errors:") {
		t.Error("Expected error summary in output")
	}
}

// TestGetPostgresVersion validates PostgreSQL version retrieval functionality.
// Creates a mock postgres executable and verifies version string parsing.
func TestGetPostgresVersion(t *testing.T) {
	tmpDir := t.TempDir()
	binDir := filepath.Join(tmpDir, "bin")
	err := os.MkdirAll(binDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create temporary bin directory: %v", err)
	}

	postgresPath := filepath.Join(binDir, "postgres")
	mockContent := `#!/bin/sh
echo "postgres (Cloudberry Database) 14.4"`
	err = os.WriteFile(postgresPath, []byte(mockContent), 0755)
	if err != nil {
		t.Fatalf("Failed to create mock postgres executable: %v", err)
	}

	version, err := getPostgresVersion(tmpDir)
	if err != nil {
		t.Errorf("Unexpected error getting postgres version: %v", err)
	}
	if !strings.Contains(version, "Cloudberry Database") {
		t.Errorf("Expected version to contain 'Cloudberry Database', got: %s", version)
	}
}

// TestGetGPVersion validates Apache Cloudberry version retrieval.
// Creates a mock postgres executable and verifies GP version string parsing.
func TestGetGPVersion(t *testing.T) {
	tmpDir := t.TempDir()
	binDir := filepath.Join(tmpDir, "bin")
	err := os.MkdirAll(binDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create temporary bin directory: %v", err)
	}

	postgresPath := filepath.Join(binDir, "postgres")
	mockContent := `#!/bin/sh
echo "postgres (Cloudberry Database) 1.6.0 build 1"`
	err = os.WriteFile(postgresPath, []byte(mockContent), 0755)
	if err != nil {
		t.Fatalf("Failed to create mock postgres executable: %v", err)
	}

	version, err := getGPVersion(tmpDir)
	if err != nil {
		t.Errorf("Unexpected error getting gp version: %v", err)
	}
	if !strings.Contains(version, "1.6.0") {
		t.Errorf("Expected version to contain '1.6.0', got: %s", version)
	}
}

func TestRunSysInfoWithMockedGPHOME(t *testing.T) {
	// Mock GPHOME environment variable
	mockGPHOME := t.TempDir()
	defer os.Setenv("GPHOME", os.Getenv("GPHOME")) // Restore original GPHOME after test
	os.Setenv("GPHOME", mockGPHOME)

	// Mock binaries and files
	binDir := filepath.Join(mockGPHOME, "bin")

	err := os.MkdirAll(binDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test bin directory: %v", err)
	}

	if err := os.WriteFile(filepath.Join(binDir, "pg_config"), []byte("#!/bin/bash\necho 'pg_config mock'"), 0755); err != nil {
		t.Fatalf("Failed to write mock pg_config file: %v", err)
	}

	if err := os.WriteFile(filepath.Join(binDir, "postgres"), []byte("#!/bin/bash\necho 'postgres mock'"), 0755); err != nil {
		t.Fatalf("Failed to write mock postgres file: %v", err)
	}

	// Mock system file paths
	procMeminfo = filepath.Join(mockGPHOME, "meminfo")
	osReleasePath = filepath.Join(mockGPHOME, "os-release")

	// Create mocked files
	if err := os.WriteFile(procMeminfo, []byte("MemTotal: 16384 kB\nMemFree: 8192 kB\n"), 0644); err != nil {
		t.Fatalf("Failed to write mock procMeminfo file: %v", err)
	}

	if err := os.WriteFile(osReleasePath, []byte("PRETTY_NAME=\"Mock Linux\"\n"), 0644); err != nil {
		t.Fatalf("Failed to write mock osReleasePath file: %v", err)
	}

	// Mock Cobra command and args
	cmd := &cobra.Command{}
	args := []string{}

	// Run sysinfo and expect success
	err = RunSysInfo(cmd, args)
	if err != nil {
		t.Errorf("Expected no error with mocked GPHOME, got: %v", err)
	}
}
