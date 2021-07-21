package handler

import (
	"github.com/Sho2010/k8s-job-notifier/internal/event"
)

type Handler interface {
	Handle(e event.Event)
}

