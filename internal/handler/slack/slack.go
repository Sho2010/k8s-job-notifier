package slack

import (
	"fmt"
	"log"
	"strings"

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
	Token            string
	DefaultChannel   string
	Title            string
	NotifyCondisions []string
}

const (
	annotationPrefix = "notify.sho2010.dev/"

	// ChannelAnnotation is annotation key of slack message destination channel
	//
	// const: 'notify.sho2010.dev/channel'
	ChannelAnnotation = annotationPrefix + "channel"

	// NotifyConditionAnnotation is
	//
	// value example: "Complete,Failed"
	NotifyConditionAnnotation = annotationPrefix + "condition"
)

// Handle handles the notification.
func (s *Slack) Handle(e event.Event) {

	//TODO: おそらく起動したときはConditions == 0 で判定できるはず
	// job createのときはconditionsが空でくる、他にいい判定方法があればそれに変える
	job := e.Resource.(*batchv1.Job)
	if len(job.Status.Conditions) == 0 {
		return
	}

	annotations := job.GetAnnotations()

	notifyCondisions := strings.Split(",", annotations[ChannelAnnotation])
	isSend := false
	for _, con := range notifyCondisions {
		con = strings.TrimSpace(con)
		con = strings.ToLower(con)

		if con == "" {
			isSend = true
			break
		}
	}

	if !isSend {
		log.Printf("ignore\n")
		return
	}
	channel := annotations[ChannelAnnotation]
	if len(channel) == 0 {
		channel = s.DefaultChannel
	}

	attachment := buildAttachment(e, s)

	client := slack.New(s.Token)
	channelID, timestamp, err := client.PostMessage(channel,
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionAsUser(false),
		slack.MsgOptionIconEmoji(":sushi:"))
	if err != nil {
		log.Printf("slack error: %s\n", err)
		return
	}

	log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
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
				Title: s.Title,
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
