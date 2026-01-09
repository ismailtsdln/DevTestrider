package report

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"

	"github.com/ismailtsdln/DevTestrider/internal/engine"
)

func GeneratePDF(result *engine.TestResult, outputDir string) (string, error) {
	if outputDir == "" {
		outputDir = "."
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", err
	}

	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(10, 15, 10)

	// Header
	m.Row(20, func() {
		m.Col(12, func() {
			m.Text("DevTestrider Report", props.Text{
				Size:  18,
				Style: consts.Bold,
				Align: consts.Center,
			})
		})
	})
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text(fmt.Sprintf("Generated: %s", result.Timestamp.Format("Jan 02, 2006 15:04")), props.Text{
				Size:  10,
				Align: consts.Center,
				Color: color.Color{Red: 100, Green: 100, Blue: 100},
			})
		})
	})

	// Status Line
	status := "PASSED"
	statusColor := color.Color{Red: 0, Green: 150, Blue: 0}
	if !result.Success {
		status = "FAILED"
		statusColor = color.Color{Red: 200, Green: 0, Blue: 0}
	}

	m.Row(20, func() {
		m.Col(12, func() {
			m.Text(status, props.Text{
				Size:  16,
				Style: consts.Bold,
				Align: consts.Center,
				Color: statusColor,
			})
		})
	})

	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Summary", props.Text{Style: consts.Bold, Size: 12})
		})
	})

	// Stats Table
	m.Row(15, func() {
		m.Col(3, func() { m.Text(fmt.Sprintf("Total: %d", result.TotalTests), props.Text{}) })
		m.Col(3, func() { m.Text(fmt.Sprintf("Passed: %d", result.PassedTests), props.Text{}) })
		m.Col(3, func() { m.Text(fmt.Sprintf("Failed: %d", result.FailedTests), props.Text{}) })
		m.Col(3, func() { m.Text(fmt.Sprintf("Duration: %.2fs", result.Duration), props.Text{}) })
	})

	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Package Details", props.Text{Style: consts.Bold, Size: 12, Top: 5})
		})
	})

	// Header Row
	m.Row(10, func() {
		m.Col(6, func() { m.Text("Package", props.Text{Style: consts.Bold}) })
		m.Col(2, func() { m.Text("Status", props.Text{Style: consts.Bold}) })
		m.Col(2, func() { m.Text("Coverage", props.Text{Style: consts.Bold}) })
		m.Col(2, func() { m.Text("Duration", props.Text{Style: consts.Bold}) })
	})

	for _, pkg := range result.Packages {
		pkgName := pkg.Name
		pkgStatus := "FAIL"
		if pkg.Status == "PASS" {
			pkgStatus = "PASS"
		}

		cov := "-"
		if pkg.Coverage > 0 {
			cov = fmt.Sprintf("%.1f%%", pkg.Coverage) // Fixed undefined 'm' error by using string formatting here
		}

		duration := fmt.Sprintf("%.3fs", pkg.Duration)

		m.Row(8, func() {
			m.Col(6, func() { m.Text(pkgName, props.Text{Size: 9}) })
			m.Col(2, func() { m.Text(pkgStatus, props.Text{Size: 9}) })
			m.Col(2, func() { m.Text(cov, props.Text{Size: 9}) })
			m.Col(2, func() { m.Text(duration, props.Text{Size: 9}) })
		})
	}

	filename := filepath.Join(outputDir, fmt.Sprintf("report-%d.pdf", time.Now().Unix()))
	if err := m.OutputFileAndClose(filename); err != nil {
		return "", err
	}

	return filename, nil
}
