package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Test case structure for JSON serialization
type TestCase struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Category        string   `json:"category"`
	FileA           string   `json:"file_a,omitempty"`
	FileB           string   `json:"file_b,omitempty"`
	FileDiff        string   `json:"file_diff,omitempty"`
	ContentA        string   `json:"content_a,omitempty"`
	ContentB        string   `json:"content_b,omitempty"`
	ExpectedDiff    string   `json:"expected_diff,omitempty"`
	Args            []string `json:"args,omitempty"`
	ExpectedExit    int      `json:"expected_exit"`
	ShouldError     bool     `json:"should_error"`
	ComplianceLevel string   `json:"compliance_level"` // "core", "extended", "format"
}

// Test results
type TestResult struct {
	TestCase   TestCase `json:"test_case"`
	Passed     bool     `json:"passed"`
	ActualExit int      `json:"actual_exit"`
	ActualDiff string   `json:"actual_diff,omitempty"`
	Error      string   `json:"error,omitempty"`
	Duration   string   `json:"duration"`
}

// Test runner configuration
type Config struct {
	BinaryPath     string
	TestDataDir    string
	Verbose        bool
	CoreOnly       bool
	ExtendedOnly   bool
	FormatOnly     bool
	CategoryFilter string
	Timeout        time.Duration
	FailFast       bool
}

func main() {
	var config Config

	flag.StringVar(&config.BinaryPath, "binary", "", "Path to jd binary to test (required)")
	flag.StringVar(&config.TestDataDir, "testdata", "./testdata", "Path to test data directory")
	flag.BoolVar(&config.Verbose, "verbose", false, "Verbose output")
	flag.BoolVar(&config.CoreOnly, "core-only", false, "Test core compliance only")
	flag.BoolVar(&config.ExtendedOnly, "extended-only", false, "Test extended compliance only")
	flag.BoolVar(&config.FormatOnly, "format-only", false, "Test format compliance only")
	flag.StringVar(&config.CategoryFilter, "category", "", "Filter by test category")
	flag.DurationVar(&config.Timeout, "timeout", 30*time.Second, "Timeout per test case")
	flag.BoolVar(&config.FailFast, "fail-fast", false, "Stop on first failure")

	flag.Parse()

	if config.BinaryPath == "" {
		if len(flag.Args()) > 0 {
			config.BinaryPath = flag.Args()[0]
		} else {
			fmt.Fprintf(os.Stderr, "Usage: %s [options] <binary-path>\n", os.Args[0])
			flag.PrintDefaults()
			os.Exit(64) // EX_USAGE
		}
	}

	// Validate binary exists and is executable
	if err := validateBinary(config.BinaryPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(66) // EX_NOINPUT
	}

	// Load test cases
	testCases, err := loadTestCases(config.TestDataDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading test cases: %v\n", err)
		os.Exit(65) // EX_DATAERR
	}

	// Filter test cases
	testCases = filterTestCases(testCases, config)

	if len(testCases) == 0 {
		fmt.Println("No test cases match the specified criteria")
		os.Exit(0)
	}

	fmt.Printf("Running %d test cases against %s\n", len(testCases), config.BinaryPath)
	fmt.Println()

	// Run tests
	results := make([]TestResult, 0, len(testCases))
	passed := 0
	failed := 0

	for i, testCase := range testCases {
		if config.Verbose {
			fmt.Printf("[%d/%d] %s: %s\n", i+1, len(testCases), testCase.Category, testCase.Name)
		}

		result := runTestCase(testCase, config)
		results = append(results, result)

		if result.Passed {
			passed++
			if config.Verbose {
				fmt.Printf("  ✓ PASS (%s)\n", result.Duration)
			}
		} else {
			failed++
			if config.Verbose || !config.FailFast {
				fmt.Printf("  ✗ FAIL (%s): %s\n", result.Duration, result.Error)
			}
			if config.FailFast {
				break
			}
		}
	}

	// Print summary
	fmt.Println()
	fmt.Printf("Results: %d passed, %d failed, %d total\n", passed, failed, len(results))

	if failed > 0 {
		fmt.Println("\nFailures:")
		for _, result := range results {
			if !result.Passed {
				fmt.Printf("  %s/%s: %s\n", result.TestCase.Category, result.TestCase.Name, result.Error)
			}
		}
	}

	// Generate detailed report if requested
	if config.Verbose {
		reportFile := "test-results.json"
		if err := generateReport(results, reportFile); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not generate report: %v\n", err)
		} else {
			fmt.Printf("Detailed report written to %s\n", reportFile)
		}
	}

	// Exit with appropriate code
	if failed > 0 {
		os.Exit(1) // Test failures
	}
	os.Exit(0) // All tests passed
}

func validateBinary(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("cannot access binary: %v", err)
	}

	if info.IsDir() {
		return fmt.Errorf("path is a directory, not a binary")
	}

	// Check if executable (Unix-style check)
	if info.Mode()&0111 == 0 {
		return fmt.Errorf("binary is not executable")
	}

	return nil
}

func loadTestCases(testDataDir string) ([]TestCase, error) {
	testCases := make([]TestCase, 0)

	// Load from cases directory
	casesDir := filepath.Join(testDataDir, "../cases")

	err := filepath.Walk(casesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".json") {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		var cases []TestCase
		if err := json.Unmarshal(data, &cases); err != nil {
			return fmt.Errorf("error parsing %s: %v", path, err)
		}

		testCases = append(testCases, cases...)
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Validate and resolve file paths
	for i := range testCases {
		if err := resolveTestCase(&testCases[i], testDataDir); err != nil {
			return nil, fmt.Errorf("error in test case %s: %v", testCases[i].Name, err)
		}
	}

	return testCases, nil
}

func resolveTestCase(testCase *TestCase, testDataDir string) error {
	// Resolve file paths relative to testdata directory
	if testCase.FileA != "" {
		testCase.FileA = filepath.Join(testDataDir, testCase.FileA)
	}
	if testCase.FileB != "" {
		testCase.FileB = filepath.Join(testDataDir, testCase.FileB)
	}
	if testCase.FileDiff != "" {
		testCase.FileDiff = filepath.Join(testDataDir, testCase.FileDiff)
	}

	// Set default compliance level
	if testCase.ComplianceLevel == "" {
		testCase.ComplianceLevel = "core"
	}

	return nil
}

func filterTestCases(testCases []TestCase, config Config) []TestCase {
	filtered := make([]TestCase, 0)

	for _, testCase := range testCases {
		// Filter by compliance level
		if config.CoreOnly && testCase.ComplianceLevel != "core" {
			continue
		}
		if config.ExtendedOnly && testCase.ComplianceLevel != "extended" {
			continue
		}
		if config.FormatOnly && testCase.ComplianceLevel != "format" {
			continue
		}

		// Filter by category
		if config.CategoryFilter != "" && testCase.Category != config.CategoryFilter {
			continue
		}

		filtered = append(filtered, testCase)
	}

	return filtered
}

func runTestCase(testCase TestCase, config Config) TestResult {
	start := time.Now()
	result := TestResult{
		TestCase: testCase,
		Passed:   false,
	}

	// Build command arguments
	args := make([]string, 0)
	args = append(args, testCase.Args...)

	// Add input files or content
	var inputFiles []string
	if testCase.ContentA != "" || (testCase.ContentA == "" && testCase.FileA == "") {
		// Create temporary file for content A (including empty content)
		fileA, err := createTempFile(testCase.ContentA)
		if err != nil {
			result.Error = fmt.Sprintf("failed to create temp file A: %v", err)
			result.Duration = time.Since(start).String()
			return result
		}
		defer os.Remove(fileA)
		inputFiles = append(inputFiles, fileA)
	} else if testCase.FileA != "" {
		inputFiles = append(inputFiles, testCase.FileA)
	}

	if testCase.ContentB != "" || (testCase.ContentB == "" && testCase.FileB == "") {
		// Create temporary file for content B (including empty content)
		fileB, err := createTempFile(testCase.ContentB)
		if err != nil {
			result.Error = fmt.Sprintf("failed to create temp file B: %v", err)
			result.Duration = time.Since(start).String()
			return result
		}
		defer os.Remove(fileB)
		inputFiles = append(inputFiles, fileB)
	} else if testCase.FileB != "" {
		inputFiles = append(inputFiles, testCase.FileB)
	}

	args = append(args, inputFiles...)

	// Execute command with timeout
	cmd := exec.Command(config.BinaryPath, args...)

	output, err := executeWithTimeout(cmd, config.Timeout)
	result.Duration = time.Since(start).String()

	// Handle execution results
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			result.ActualExit = exitError.ExitCode()
		} else {
			result.Error = fmt.Sprintf("execution error: %v", err)
			return result
		}
	} else {
		result.ActualExit = 0
	}

	result.ActualDiff = string(output)

	// Check exit code
	if result.ActualExit != testCase.ExpectedExit {
		result.Error = fmt.Sprintf("exit code mismatch: expected %d, got %d",
			testCase.ExpectedExit, result.ActualExit)
		return result
	}

	// Check for expected errors
	if testCase.ShouldError && result.ActualExit == 0 {
		result.Error = "expected error but command succeeded"
		return result
	}

	// Compare output if provided
	if testCase.ExpectedDiff != "" {
		if !compareOutputs(result.ActualDiff, testCase.ExpectedDiff) {
			result.Error = fmt.Sprintf("output mismatch:\nExpected:\n%s\nActual:\n%s",
				testCase.ExpectedDiff, result.ActualDiff)
			return result
		}
	}

	result.Passed = true
	return result
}

func createTempFile(content string) (string, error) {
	tmpFile, err := ioutil.TempFile("", "jd-test-*.json")
	if err != nil {
		return "", err
	}

	_, err = tmpFile.WriteString(content)
	if err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", err
	}

	tmpFile.Close()
	return tmpFile.Name(), nil
}

func executeWithTimeout(cmd *exec.Cmd, timeout time.Duration) ([]byte, error) {
	done := make(chan error, 1)
	var output []byte
	var err error

	go func() {
		output, err = cmd.Output()
		done <- err
	}()

	select {
	case err := <-done:
		return output, err
	case <-time.After(timeout):
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return nil, fmt.Errorf("command timed out after %v", timeout)
	}
}

func compareOutputs(actual, expected string) bool {
	// Normalize line endings
	actual = strings.ReplaceAll(actual, "\r\n", "\n")
	expected = strings.ReplaceAll(expected, "\r\n", "\n")

	// Trim trailing whitespace
	actual = strings.TrimSpace(actual)
	expected = strings.TrimSpace(expected)

	// Collect expected metadata lines so we know which ones are required
	expectedMeta := make(map[string]bool)
	for _, line := range strings.Split(expected, "\n") {
		if strings.HasPrefix(line, "^ ") {
			expectedMeta[line] = true
		}
	}

	// Filter actual output: drop metadata lines not present in expected
	var filteredLines []string
	for _, line := range strings.Split(actual, "\n") {
		if strings.HasPrefix(line, "^ ") && !expectedMeta[line] {
			continue
		}
		filteredLines = append(filteredLines, line)
	}
	actual = strings.Join(filteredLines, "\n")

	return actual == expected
}

func generateReport(results []TestResult, filename string) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}
