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

// File: cmd/core_parser_output.go
// Purpose: Implements utilities for saving and comparing core dump analysis results.
// Includes functions to serialize analysis data in JSON or YAML format and
// to identify common patterns across multiple core files.
// Dependencies: Uses Go standard libraries for JSON, YAML, file handling, and string manipulation.

package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

// saveAnalysis saves analysis results to a file.
// Parameters:
// - analysis: The CoreAnalysis object containing the data to save.
// Returns:
// - An error if the operation fails, or nil otherwise.
func saveAnalysis(analysis CoreAnalysis) error {
	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(outputDir, fmt.Sprintf("core_analysis_%s.%s", timestamp, formatFlag))

	// Deduplicate threads and parse basic info
	analysis.Threads = deduplicateThreads(analysis.Threads)
	analysis.BasicInfo = parseBasicInfo(analysis.FileInfo.FileOutput)

	// Mark crashed threads and enhance thread info
	for i := range analysis.Threads {
		for _, frame := range analysis.Threads[i].Backtrace {
			if strings.Contains(frame.Function, "SigillSigsegvSigbus") {
				analysis.Threads[i].IsCrashed = true
				break
			}
		}
		analysis.Threads[i].Name = determineThreadRole(analysis.Threads[i].Backtrace)
	}

	var data []byte
	var err error
	if formatFlag == "json" {
		data, err = json.MarshalIndent(analysis, "", "  ")
	} else {
		data, err = yaml.Marshal(analysis)
	}
	if err != nil {
		return fmt.Errorf("failed to marshal analysis: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write analysis file: %w", err)
	}

	fmt.Printf("Analysis saved to: %s\n", filename)
	return nil
}

// compareCores analyzes multiple core files to identify patterns.
// Parameters:
// - analyses: A slice of CoreAnalysis objects representing individual core dump analyses.
// Returns:
// - A CoreComparison object summarizing common patterns and statistics.
func compareCores(analyses []CoreAnalysis) CoreComparison {
	comparison := CoreComparison{
		TotalCores:      len(analyses),
		CommonSignals:   make(map[string]int),
		CommonFunctions: make(map[string]int),
		TimeRange:       make(map[string]string),
	}

	// Track time range
	var firstTime, lastTime time.Time
	for i, analysis := range analyses {
		t, _ := time.Parse(time.RFC3339, analysis.Timestamp)
		if i == 0 || t.Before(firstTime) {
			firstTime = t
		}
		if i == 0 || t.After(lastTime) {
			lastTime = t
		}
	}
	comparison.TimeRange["first"] = firstTime.Format(time.RFC3339)
	comparison.TimeRange["last"] = lastTime.Format(time.RFC3339)

	// Collect signal and function distributions
	crashGroups := make(map[string][]CoreAnalysis)
	for _, analysis := range analyses {
		signal := analysis.SignalInfo.SignalName
		comparison.CommonSignals[signal]++

		// Count functions in stack traces
		for _, frame := range analysis.StackTrace {
			if !isSystemFunction(frame.Function) {
				comparison.CommonFunctions[frame.Function]++
			}
		}

		// Create crash signature
		var signature strings.Builder
		signature.WriteString(signal)
		for i, frame := range analysis.StackTrace {
			if i < 3 && !isSystemFunction(frame.Function) {
				signature.WriteString("|" + frame.Function)
			}
		}
		crashGroups[signature.String()] = append(crashGroups[signature.String()], analysis)
	}

	// Generate crash patterns
	for signature, group := range crashGroups {
		if len(group) > 1 {
			parts := strings.Split(signature, "|")
			pattern := CrashPattern{
				Signal:            parts[0],
				StackSignature:    parts[1:],
				OccurrenceCount:   len(group),
				AffectedCoreFiles: make([]string, 0, len(group)),
			}
			for _, analysis := range group {
				pattern.AffectedCoreFiles = append(pattern.AffectedCoreFiles, analysis.CoreFile)
			}
			comparison.CrashPatterns = append(comparison.CrashPatterns, pattern)
		}
	}

	// Sort patterns by occurrence count
	sort.Slice(comparison.CrashPatterns, func(i, j int) bool {
		return comparison.CrashPatterns[i].OccurrenceCount > comparison.CrashPatterns[j].OccurrenceCount
	})

	return comparison
}

// saveComparison saves comparison results to a file.
// Parameters:
// - comparison: The CoreComparison object summarizing core file patterns.
// Returns:
// - An error if the operation fails, or nil otherwise.
func saveComparison(comparison CoreComparison) error {
	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(outputDir, fmt.Sprintf("core_comparison_%s.%s", timestamp, formatFlag))

	var data []byte
	var err error
	if formatFlag == "json" {
		data, err = json.MarshalIndent(comparison, "", "  ")
	} else {
		data, err = yaml.Marshal(comparison)
	}
	if err != nil {
		return fmt.Errorf("failed to marshal comparison: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write comparison file: %w", err)
	}

	fmt.Printf("Comparison results saved to: %s\n", filename)
	return nil
}
