package orchestrator

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/ismailtsdln/DevTestrider/internal/config"
	"github.com/ismailtsdln/DevTestrider/internal/engine"
	"github.com/ismailtsdln/DevTestrider/internal/notify"
	"github.com/ismailtsdln/DevTestrider/internal/report"
	"github.com/ismailtsdln/DevTestrider/internal/server"
)

type Orchestrator struct {
	cfg     *config.Config
	runner  *engine.Runner
	watcher *engine.Watcher
	server  *server.Server
}

func New(cfg *config.Config, r *engine.Runner, w *engine.Watcher, s *server.Server) *Orchestrator {
	return &Orchestrator{
		cfg:     cfg,
		runner:  r,
		watcher: w,
		server:  s,
	}
}

func (o *Orchestrator) Start(done chan bool) {
	// Styles
	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	passStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	failStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("160"))

	for {
		select {
		case eventPath := <-o.watcher.Events:
			fmt.Printf("\n%s %s\n", infoStyle.Render("File changed:"), eventPath)

			// Run Tests
			result, err := o.runner.RunTests(filepath.Dir(eventPath))
			if err != nil {
				log.Printf("Error running tests: %v", err)
				continue
			}

			// Render Status Header
			status := failStyle.Render("FAILED ❌")
			if result.Success {
				status = passStyle.Render("PASSED ✅")
			}
			fmt.Printf("Status: %s (Duration: %.2fs)\n", status, result.Duration)

			// Render Package Details
			for _, pkg := range result.Packages {
				pkgStatus := failStyle.Render("FAIL")
				if pkg.Status == "PASS" {
					pkgStatus = passStyle.Render("PASS")
				}

				covStr := "N/A"
				if pkg.Coverage > 0 {
					covColor := "160" // Red
					if pkg.Coverage > 50 {
						covColor = "220"
					} // Yellow
					if pkg.Coverage > 80 {
						covColor = "42"
					} // Green
					covStr = lipgloss.NewStyle().Foreground(lipgloss.Color(covColor)).Render(fmt.Sprintf("%.1f%%", pkg.Coverage))
				}

				fmt.Printf("  • %-40s %s  %s (%.2fs)\n",
					pkg.Name,
					pkgStatus,
					covStr,
					pkg.Duration,
				)
			}

			// Generate Reports
			if len(o.cfg.Report.Formats) > 0 {
				for _, fmtType := range o.cfg.Report.Formats {
					var path string
					var err error
					switch fmtType {
					case "html":
						path, err = report.GenerateHTML(result, o.cfg.Report.OutputDir)
					case "pdf":
						path, err = report.GeneratePDF(result, o.cfg.Report.OutputDir)
					}

					if err != nil {
						log.Printf("Failed to generate %s report: %v", fmtType, err)
					} else if path != "" {
						fmt.Printf("Report generated: %s\n", lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Render(path))
					}
				}
			}

			// Notifications & Broadcast
			notify.SendNotification(o.cfg.Notifications, result)
			o.server.Broadcast(result)

		case <-done:
			return
		}
	}
}
