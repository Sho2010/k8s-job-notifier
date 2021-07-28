package handler

import (
	"os"

	"github.com/Sho2010/k8s-job-notifier/internal/handler/slack"
)

func CreateHandler() (Handler, error) {
	h, err := CreateSlackHandler()

	if err != nil {
		return nil, err
	}

	return h, nil
}

func CreateSlackHandler() (Handler, error) {
	dc := os.Getenv("DEFAULT_CHANNEL")
	if len(dc) == 0 {
		dc = "#bot_sandbox"
	}

	return &slack.Slack{
		Token:            os.Getenv("SLACK_TOKEN"),
		DefaultChannel:   dc,
		Title:            "job notify",
		NotifyCondisions: []string{"Failed"},
	}, nil
}
