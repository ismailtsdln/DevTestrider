package report

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/props"

	"github.com/ismailtsdln/DevTestrider/internal/engine"
)

func GeneratePDF(result *engine.TestResult, outputDir string) (string, error) {
	if outputDir == "" {
		outputDir = "."
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", err
	}

	cfg := config.NewBuilder().
		WithPageSize(config.A4).
		WithMargins(10, 15, 10).
		Build()

	m := maroto.New(cfg)

	// Header
	m.AddRow(20,
		text.NewCol(12, "DevTestrider Report",
			props.Text{
				Size:  18,
				Style: fontstyle.Bold,
				Align: align.Center,
			}),
	)
	m.AddRow(10,
		text.NewCol(12, fmt.Sprintf("Generated: %s", result.Timestamp.Format("Jan 02, 2006 15:04")),
			props.Text{
				Size:  10,
				Align: align.Center,
				Color: &props.Color{Red: 100, Green: 100, Blue: 100},
			}),
	)

	// Status Line
	status := "PASSED"
	statusColor := &props.Color{Red: 0, Green: 150, Blue: 0}
	if !result.Success {
		status = "FAILED"
		statusColor = &props.Color{Red: 200, Green: 0, Blue: 0}
	}

	m.AddRow(20,
		text.NewCol(12, status, props.Text{
			Size:  16,
			Style: fontstyle.Bold,
			Align: align.Center,
			Color: statusColor,
		}),
	)

	m.AddRow(10, text.NewCol(12, "Summary", props.Text{Style: fontstyle.Bold, Size: 12}))

	// Stats Table
	m.AddRow(15,
		col.New(3).Add(text.New(fmt.Sprintf("Total: %d", result.TotalTests))),
		col.New(3).Add(text.New(fmt.Sprintf("Passed: %d", result.PassedTests))),
		col.New(3).Add(text.New(fmt.Sprintf("Failed: %d", result.FailedTests))),
		col.New(3).Add(text.New(fmt.Sprintf("Duration: %.2fs", result.Duration))),
	)

	m.AddRow(10, text.NewCol(12, "Package Details", props.Text{Style: fontstyle.Bold, Size: 12, Top: 5}))

	// Header Row
	m.AddRow(10,
		text.NewCol(6, "Package", props.Text{Style: fontstyle.Bold}),
		text.NewCol(2, "Status", props.Text{Style: fontstyle.Bold}),
		text.NewCol(2, "Coverage", props.Text{Style: fontstyle.Bold}),
		text.NewCol(2, "Duration", props.Text{Style: fontstyle.Bold}),
	)

	for _, pkg := range result.Packages {
		pkgStatus := "FAIL"
		if pkg.Status == "PASS" {
			pkgStatus = "PASS"
		}

		cov := "-"
		if pkg.Coverage > 0 {
			cov = fmt.Sprintf("%.1f%%", pkg.Coverage)
		}

		m.AddRow(8,
			text.NewCol(6, pkg.Name, props.Text{Size: 9}),
			text.NewCol(2, pkgStatus, props.Text{Size: 9}),
			text.NewCol(2, cov, props.Text{Size: 9}),
			text.NewCol(2, fmt.Sprintf("%.3fs", pkg.Duration), props.Text{Size: 9}),
		)
	}

	filename := filepath.Join(outputDir, fmt.Sprintf("report-%d.pdf", time.Now().Unix()))
	document, err := m.Generate()
	if err != nil {
		return "", err
	}

	if err := document.Save(filename); err != nil {
		return "", err
	}

	return filename, nil
}
