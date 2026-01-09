package report

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/ismailtsdln/DevTestrider/internal/engine"
)

const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>DevTestrider Report</title>
    <style>
        body { font-family: system-ui, -apple-system, sans-serif; background: #0f172a; color: #f8fafc; padding: 2rem; }
        .container { max-width: 1000px; margin: 0 auto; }
        .header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; padding-bottom: 1rem; border-bottom: 1px solid #1e293b; }
        .success { color: #34d399; }
        .failure { color: #f43f5e; }
        .card { background: #1e293b; border-radius: 0.5rem; padding: 1.5rem; margin-bottom: 1rem; }
        .stats { display: grid; grid-template-columns: repeat(4, 1fr); gap: 1rem; margin-bottom: 2rem; }
        .stat-box { background: #1e293b; padding: 1rem; border-radius: 0.5rem; text-align: center; }
        .stat-value { font-size: 1.5rem; font-weight: bold; display: block; }
        .stat-label { font-size: 0.875rem; color: #94a3b8; }
        table { width: 100%; border-collapse: collapse; }
        th, td { text-align: left; padding: 1rem; border-bottom: 1px solid #334155; }
        th { color: #94a3b8; font-weight: 500; font-size: 0.875rem; }
        .badge { padding: 0.25rem 0.5rem; border-radius: 0.25rem; font-size: 0.75rem; font-weight: 600; }
        .badge-pass { background: rgba(52, 211, 153, 0.1); color: #34d399; }
        .badge-fail { background: rgba(244, 63, 94, 0.1); color: #f43f5e; }
    </style>
</head>
<body>
    <div class="container">
        <header class="header">
            <div>
                <h1 style="margin:0">DevTestrider Report</h1>
                <p style="color: #64748b; margin: 0.5rem 0 0;">Generated on {{.Timestamp.Format "Jan 02, 2006 15:04:05"}}</p>
            </div>
            <div>
                {{if .Success}}
                    <span class="success" style="font-size: 1.25rem; font-weight: bold;">PASSED ✅</span>
                {{else}}
                    <span class="failure" style="font-size: 1.25rem; font-weight: bold;">FAILED ❌</span>
                {{end}}
            </div>
        </header>

        <div class="stats">
            <div class="stat-box">
                <span class="stat-value">{{.TotalTests}}</span>
                <span class="stat-label">Total Tests</span>
            </div>
            <div class="stat-box">
                <span class="stat-value success">{{.PassedTests}}</span>
                <span class="stat-label">Passed</span>
            </div>
            <div class="stat-box">
                <span class="stat-value failure">{{.FailedTests}}</span>
                <span class="stat-label">Failed</span>
            </div>
            <div class="stat-box">
                <span class="stat-value">{{printf "%.2f" .Duration}}s</span>
                <span class="stat-label">Duration</span>
            </div>
        </div>

        <div class="card">
            <h3>Package Results</h3>
            <table>
                <thead>
                    <tr>
                        <th>Package</th>
                        <th>Tests</th>
                        <th>Coverage</th>
                        <th>Duration</th>
                        <th>Status</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .Packages}}
                    <tr>
                        <td>{{.Name}}</td>
                        <td>{{len .Tests}}</td>
                        <td>{{if gt .Coverage 0.0}}{{printf "%.1f%%" .Coverage}}{{else}}-{{end}}</td>
                        <td>{{printf "%.3f" .Duration}}s</td>
                        <td>
                            {{if eq .Status "PASS"}}
                                <span class="badge badge-pass">PASS</span>
                            {{else}}
                                <span class="badge badge-fail">FAIL</span>
                            {{end}}
                        </td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
        </div>
    </div>
</body>
</html>
`

func GenerateHTML(result *engine.TestResult, outputDir string) (string, error) {
	if outputDir == "" {
		outputDir = "."
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", err
	}

	filename := filepath.Join(outputDir, fmt.Sprintf("report-%d.html", time.Now().Unix()))
	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	tmpl, err := template.New("report").Parse(htmlTemplate)
	if err != nil {
		return "", err
	}

	if err := tmpl.Execute(f, result); err != nil {
		return "", err
	}

	return filename, nil
}
