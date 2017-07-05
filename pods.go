package main

import (
	"log"
	"time"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

func newPodHandler(obj interface{}) {
	pod, ok := obj.(*v1.Pod)
	if !ok {
		log.Printf("unexpected object type: %#v\n", obj)
	}
	log.Println("pod-event> add:", pod.GetName())

	podmetrics.TotalNumberPods.Add(1)
}

func deletePodHandler(obj interface{}) {
	pod, ok := obj.(*v1.Pod)
	if !ok {
		log.Printf("unexpected object type: %#v\n", obj)
	}
	log.Println("pod-event> delete:", pod.GetName())

	podmetrics.TotalNumberPods.Sub(1)
	// TODO ? delete metrics for the pod
}

func RunPodInformer(stop chan struct{}) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}

	cs, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}

	// query list of current pods
	// podlist, err := cs.Core().Pods(v1.NamespaceDefault).List(metav1.ListOptions{})
	// if err != nil {
	// 	log.Printf("Error listing pods: %s", err)
	// }

	// fmt.Println("Found existing pods:")
	// for _, pod := range podlist.Items {
	// 	fmt.Printf("  > %s\n", pod.GetName())
	// }

	// watch new pod events
	lw := cache.NewListWatchFromClient(cs.Core().RESTClient(), "pods", v1.NamespaceDefault, fields.Everything())
	_, ctl := cache.NewInformer(
		lw,
		&v1.Pod{},
		0*time.Second,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    newPodHandler,
			DeleteFunc: deletePodHandler,
		},
	)

	log.Println("started watching pod events")
	ctl.Run(stop)
}
