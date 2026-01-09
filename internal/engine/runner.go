package engine

import (
	"bufio"
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Runner struct {
	Running bool
}

func NewRunner() *Runner {
	return &Runner{}
}

func (r *Runner) RunTests(path string) (*TestResult, error) {
	r.Running = true
	defer func() { r.Running = false }()

	// If path is a file, get directory
	if strings.HasSuffix(path, ".go") {
		// We usually want to run tests for the whole package or project even if one file changed,
		// to catch regressions. For now, let's run "./..." from root or specific package.
		// A simple strategy: always run "./..." for full coverage or run specific package.
		// Let's default to ./... for now as per requirements "Autotest Engine".
		path = "./..."
	}

	cmd := exec.Command("go", "test", "-json", "-cover", path)
	cmd.Stderr = os.Stderr // Capture stderr if needed

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	result := &TestResult{
		Timestamp: time.Now(),
		Packages:  make(map[string]*PackageResult),
		Success:   true,
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Bytes()
		var event GoTestEvent
		if err := json.Unmarshal(line, &event); err != nil {
			continue // Skip non-JSON lines (e.g. build output)
		}

		r.processEvent(result, event)
	}

	if err := cmd.Wait(); err != nil {
		// go test returns exit code 1 if tests fail, which is expected
		result.Success = false
	}

	// Determine overall success if not already set by exit code (though exit code usually covers it)
	if result.FailedTests > 0 {
		result.Success = false
	}

	return result, nil
}

func (r *Runner) processEvent(result *TestResult, event GoTestEvent) {
	// Initialize package entry if needed
	if event.Package == "" {
		return
	}

	if _, exists := result.Packages[event.Package]; !exists {
		result.Packages[event.Package] = &PackageResult{
			Name:  event.Package,
			Tests: []*TestCase{},
		}
	}
	pkg := result.Packages[event.Package]

	switch event.Action {
	case "run":
		// Test started
	case "pass", "fail", "skip":
		if event.Test != "" {
			// This is a test case
			testCase := &TestCase{
				Name:     event.Test,
				Duration: event.Elapsed,
				Status:   strings.ToUpper(event.Action),
			}
			pkg.Tests = append(pkg.Tests, testCase)

			result.TotalTests++
			if event.Action == "pass" {
				result.PassedTests++
			} else if event.Action == "fail" {
				result.FailedTests++
				result.Success = false
			} else {
				result.SkippedTests++
			}
		} else {
			// This is the package result
			pkg.Status = strings.ToUpper(event.Action)
			pkg.Duration = event.Elapsed
			result.Duration += event.Elapsed
		}
	case "output":
		// Attach output to the currently running test or package log
		// Only logic for attaching to specific tests is complex without tracking state.
		// For simplicity, we won't attach granular output in this initial version
		// unless we track the 'Run' action stack.
	}
}
