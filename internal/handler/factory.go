package handler

import (
	"os"

	"github.com/Sho2010/k8s-job-notifier/internal/handler/slack"
)

func CreateHandler() (Handler, error) {
	s := slack.Slack{
		Token:   os.Getenv("WEBHOOK_URL"),
		Channel: "#bot_sandbox",
		Title:   "test",
	}

	return &s, nil
}
