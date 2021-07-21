package slack

import (
	"fmt"
	"log"

	"github.com/Sho2010/k8s-job-notifier/internal/event"
	"github.com/slack-go/slack"
	batchv1 "k8s.io/api/batch/v1"
)

var slackMsg = `
Name: %s/%s
Message: %s
Status: %s
CompletionTime: %s
`

// Slack handler implements handler.Handler interface,
// Notify event to slack channel
type Slack struct {
	Token string
	// TODO: 今動かない
	Channel string
	Title   string
}

// Handle handles the notification.
func (s *Slack) Handle(e event.Event) {
	// job createのときはconditionsが空でくる、他にいい判定方法があればそれに変える
	job := e.Resource.(*batchv1.Job)
	if len(job.Status.Conditions) == 0 {
		return
	}

	// TODO: 通知するイベントタイプを設定できるようにする
	if job.Status.Conditions[0].Type != batchv1.JobFailed {
		log.Printf("ignore other than failed\n")
		return
	}

	attachment := buildAttachment(e, s)

	// Webhook
	msg := slack.WebhookMessage{
		Attachments: []slack.Attachment{attachment},
	}
	err := slack.PostWebhook(s.Token, &msg)
	if err != nil {
		log.Printf("slack error: %s\n", err)
		return
	}

	// TODO: slack appに移行する
	// channelID, timestamp, err := client.PostMessage(s.Channel,
	// 	slack.MsgOptionAttachments(attachment),
	// 	slack.MsgOptionAsUser(true))
	// if err != nil {
	// 	log.Printf("slack error: %s\n", err)
	// 	return
	// }

	// log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	log.Printf("Message successfully sent to channel")
}

func slackColor(t batchv1.JobConditionType) string {
	switch t {
	case batchv1.JobSuspended:
		return "warning"
	case batchv1.JobComplete:
		return "good"
	case batchv1.JobFailed:
		return "danger"
	default:
		return "good"
	}
}

func buildMessage(e event.Event) string {
	job := e.Resource.(*batchv1.Job)

	// 同時実行数が1前提で作られるので仕様を考える
	s := fmt.Sprintf(slackMsg,
		job.Namespace, job.Name,
		job.Status.Conditions[0].Message,
		job.Status.Conditions[0].Type,
		job.Status.CompletionTime,
	)

	return s

}

func buildAttachment(e event.Event, s *Slack) slack.Attachment {
	mes := buildMessage(e)
	attachment := slack.Attachment{
		Fields: []slack.AttachmentField{
			{
				Title: "Job notify",
				Value: mes,
			},
		},
	}

	// TODO: とりあえずここに書くがリファクタする
	job := e.Resource.(*batchv1.Job)
	attachment.Color = slackColor(job.Status.Conditions[0].Type)
	attachment.MarkdownIn = []string{"fields"}

	return attachment
}
