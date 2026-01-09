package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/charmbracelet/lipgloss"
	"github.com/ismailtsdln/DevTestrider/internal/config"
	"github.com/ismailtsdln/DevTestrider/internal/engine"
	"github.com/ismailtsdln/DevTestrider/internal/orchestrator"
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

		// Start Watcher
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

		// Start Orchestrator
		orch := orchestrator.New(cfg, runner, watcher, srv, infoStyle, passStyle, failStyle)
		quit := make(chan os.Signal, 1)
		orchestratorDone := make(chan bool)

		go orch.Start(orchestratorDone)

		// Graceful Shutdown
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		fmt.Println(infoStyle.Render("Shutting down..."))
		// orchestratorDone <- true // Optional cleanup
		watcher.Stop()
	},
}

func Execute() error {
	return rootCmd.Execute()
}
