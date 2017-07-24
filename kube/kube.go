/*
Copyright 2017 Kinvolk GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kube

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

func WatchPodEvents(ctx context.Context, cancelFunc context.CancelFunc, cfg *rest.Config, events chan Event) {
	_, err := watchCustomResources(ctx, cfg, events)
	if err != nil {
		cancelFunc()
		return
	}

	<-ctx.Done()
}

func watchCustomResources(ctx context.Context, cfg *rest.Config, events chan Event) (cache.Controller, error) {
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	// watch new pod events
	source := cache.NewListWatchFromClient(clientset.Core().RESTClient(), string(v1.ResourcePods), v1.NamespaceDefault, fields.Everything())
	_, k8sController := cache.NewInformer(
		source,
		&v1.Pod{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    onAdd(events),
			UpdateFunc: onUpdate(events),
			DeleteFunc: onDelete(events),
		},
	)

	go k8sController.Run(ctx.Done())

	return k8sController, nil
}
