package engine

import "time"

// GoTestEvent represents a line of output from 'go test -json'
type GoTestEvent struct {
	Time    time.Time `json:"Time"`
	Action  string    `json:"Action"`
	Package string    `json:"Package"`
	Test    string    `json:"Test,omitempty"`
	Output  string    `json:"Output,omitempty"`
	Elapsed float64   `json:"Elapsed,omitempty"`
}

// TestResult represents the aggregated result of a test run
type TestResult struct {
	Timestamp    time.Time                 `json:"timestamp"`
	TotalTests   int                       `json:"total_tests"`
	PassedTests  int                       `json:"passed_tests"`
	FailedTests  int                       `json:"failed_tests"`
	SkippedTests int                       `json:"skipped_tests"`
	Duration     float64                   `json:"duration"`
	Packages     map[string]*PackageResult `json:"packages"`
	Success      bool                      `json:"success"`
}

type PackageResult struct {
	Name     string      `json:"name"`
	Duration float64     `json:"duration"`
	Status   string      `json:"status"` // PASS, FAIL
	Tests    []*TestCase `json:"tests"`
	Coverage float64     `json:"coverage"`
}

type TestCase struct {
	Name     string   `json:"name"`
	Duration float64  `json:"duration"`
	Status   string   `json:"status"` // PASS, FAIL, SKIP
	Output   []string `json:"output"`
}
