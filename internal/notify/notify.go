package notify

import (
	"fmt"
	"log"

	"github.com/gen2brain/beeep"
	"github.com/ismailtsdln/DevTestrider/internal/config"
	"github.com/ismailtsdln/DevTestrider/internal/engine"
)

func SendNotification(cfg config.NotificationsConfig, result *engine.TestResult) {
	if !cfg.Enable {
		return
	}

	title := "DevTestrider"
	message := fmt.Sprintf("Tests Passed: %d/%d", result.PassedTests, result.TotalTests)
	if !result.Success {
		message = fmt.Sprintf("Tests Failed! %d failed, %d passed", result.FailedTests, result.PassedTests)
	}

	// Desktop Notification
	if contains(cfg.Channels, "desktop") {
		err := beeep.Notify(title, message, "")
		if err != nil {
			log.Printf("Failed to send notification: %v", err)
		}
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
