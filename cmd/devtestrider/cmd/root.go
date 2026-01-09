package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/charmbracelet/lipgloss"
	"github.com/ismailtsdln/DevTestrider/internal/config"
	"github.com/ismailtsdln/DevTestrider/internal/engine"
	"github.com/ismailtsdln/DevTestrider/internal/notify"
	"github.com/ismailtsdln/DevTestrider/internal/report"
	"github.com/ismailtsdln/DevTestrider/internal/server"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "devtestrider",
	Short: "DevTestrider - Real-time Go Test Runner & Dashboard",
	Run: func(cmd *cobra.Command, args []string) {
		// Load Config
		cfgPath := "testrider.yml"
		if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
			// If not found, use defaults
			log.Println("Config file not found, using defaults")
		}

		cfg, err := config.Load(cfgPath)
		if err != nil {
			// Fallback to default if error or file missing
			// For this MVP, let's just create a default config in memory
			cfg = &config.Config{
				Watch:  config.WatchConfig{Paths: []string{"."}, Ignore: []string{".git", "node_modules", "vendor"}},
				Server: config.ServerConfig{Port: 8080},
			}
		}

		// Initialize Components
		runner := engine.NewRunner()
		srv := server.NewServer(cfg.Server)

		watcher, err := engine.NewWatcher(cfg.Watch)
		if err != nil {
			log.Fatalf("Failed to create watcher: %v", err)
		}

		// Start Server in Goroutine
		go func() {
			if err := srv.Start(); err != nil {
				log.Fatalf("Server failed: %v", err)
			}
		}()

		// Start Watcher and Logic Loop
		go watcher.Start()

		titleStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

		infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
		passStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
		failStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("160"))

		fmt.Println(titleStyle.Render("DevTestrider Started"))
		fmt.Println(infoStyle.Render("Watching for file changes..."))
		fmt.Printf("Server running at http://localhost:%d\n", cfg.Server.Port)

		// Main Loop
		go func() {
			for eventPath := range watcher.Events {
				fmt.Printf("\n%s %s\n", infoStyle.Render("File changed:"), eventPath)

				result, err := runner.RunTests(filepath.Dir(eventPath))
				if err != nil {
					log.Printf("Error running tests: %v", err)
					continue
				}

				// Header
				status := failStyle.Render("FAILED ❌")
				if result.Success {
					status = passStyle.Render("PASSED ✅")
				}
				fmt.Printf("Status: %s (Duration: %.2fs)\n", status, result.Duration)

				// Detailed Package Table
				for _, pkg := range result.Packages {
					pkgStatus := failStyle.Render("FAIL")
					if pkg.Status == "PASS" {
						pkgStatus = passStyle.Render("PASS")
					}

					// Coverage formatting
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
				if len(cfg.Report.Formats) > 0 {
					for _, fmtType := range cfg.Report.Formats {
						var path string
						var err error
						switch fmtType {
						case "html":
							path, err = report.GenerateHTML(result, cfg.Report.OutputDir)
						case "pdf":
							path, err = report.GeneratePDF(result, cfg.Report.OutputDir)
						}

						if err != nil {
							log.Printf("Failed to generate %s report: %v", fmtType, err)
						} else if path != "" {
							fmt.Printf("Report generated: %s\n", lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Render(path))
						}
					}
				}

				// Send Notification
				notify.SendNotification(cfg.Notifications, result)

				srv.Broadcast(result)
			}
		}()

		// Graceful Shutdown
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		fmt.Println(infoStyle.Render("Shutting down..."))
		watcher.Stop()
	},
}

func Execute() error {
	return rootCmd.Execute()
}
