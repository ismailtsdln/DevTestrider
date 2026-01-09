package engine

import (
	"bufio"
	"bytes"
	"os/exec"
)

type ValidationIssue struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Message string `json:"message"`
}

func RunVet(path string) ([]string, error) {
	cmd := exec.Command("go", "vet", path)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		// Vet failed (issues found or error)
		var issues []string
		scanner := bufio.NewScanner(&stderr)
		for scanner.Scan() {
			text := scanner.Text()
			if text != "" && text[0] != '#' {
				issues = append(issues, text)
			}
		}
		return issues, nil // Return issues, treat simple vet failure as "we found issues"
	}

	return []string{}, nil
}
