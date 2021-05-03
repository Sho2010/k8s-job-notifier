package controller

import (
	"context"
	"log"

	batchv1 "k8s.io/api/batch/v1"
	batchv2 "k8s.io/api/batch/v2alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type MainController struct {
	client kubernetes.Interface
}

func NewMainController(client kubernetes.Interface) MainController {
	return MainController{
		client: client,
	}
}

func createJobInformer() {

}

func (c *MainController) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobListWatcher := cache.NewListWatchFromClient(c.client.BatchV1().RESTClient(), "jobs", v1.NamespaceAll, fields.Everything())
	_, jobInformer := cache.NewIndexerInformer(jobListWatcher, &batchv1.Job{}, 0, cache.ResourceEventHandlerFuncs{

		AddFunc: func(obj interface{}) {
			c.notify(ctx, obj.(*batchv1.Job))
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			c.notify(ctx, new.(*batchv1.Job))
		},
		DeleteFunc: func(obj interface{}) {
		},
	}, cache.Indexers{})

	stopCh := make(chan struct{})
	defer close(stopCh)
	go jobInformer.Run(stopCh)
	select {} // Block all
}

func (c *MainController) notify(ctx context.Context, job *batchv1.Job) {
	// log.Printf("event: %s", job.GetName())
	log.Printf("event: %s, %s", job.GetName(), job.Status.String())
}

func (c *MainController) notify2(ctx context.Context, cj *batchv2.CronJob) {
	log.Printf("event: %s", cj.GetName())
}
