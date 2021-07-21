package main

import (
	"log"
	"os"

	"github.com/Sho2010/k8s-job-notifier/internal/controller"
	"github.com/Sho2010/k8s-job-notifier/internal/utils"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	if os.Getenv("WEBHOOK_URL") == "" {
		log.Println("WEBHOOK_URL not set")
		os.Exit(1)
	}

	var kubeClient kubernetes.Interface

	if _, err := rest.InClusterConfig(); err != nil {
		kubeClient = utils.GetClientOutOfCluster()
	} else {
		kubeClient = utils.GetClient()
	}

	c := controller.NewMainController(kubeClient)
	c.Run()
}
