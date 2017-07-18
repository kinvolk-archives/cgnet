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
	"log"
	"time"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

func RunPodInformer(stop chan struct{}, events chan Event) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}

	cs, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}

	// watch new pod events
	lw := cache.NewListWatchFromClient(cs.Core().RESTClient(), "pods", v1.NamespaceDefault, fields.Everything())
	_, ctl := cache.NewInformer(
		lw,
		&v1.Pod{},
		0*time.Second,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    emitEvent(events, NewPodEvent),
			DeleteFunc: emitEvent(events, DeletePodEvent),
		},
	)

	log.Println("started watching pod events")
	ctl.Run(stop)
}
