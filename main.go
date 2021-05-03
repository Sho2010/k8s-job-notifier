package main

import (
	"github.com/Sho2010/k8s-job-notifier/internal/controller"
	"github.com/Sho2010/k8s-job-notifier/internal/utils"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	var kubeClient kubernetes.Interface

	if _, err := rest.InClusterConfig(); err != nil {
		kubeClient = utils.GetClientOutOfCluster()
	} else {
		kubeClient = utils.GetClient()
	}

	c := controller.NewMainController(kubeClient)
	c.Run()
}
