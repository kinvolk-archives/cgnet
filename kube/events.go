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

import "k8s.io/client-go/pkg/api/v1"

type EventType int

const (
	NewPodEvent EventType = iota
	UpdatePodEvent
	DeletePodEvent
)

type Event struct {
	Type        EventType
	PodUID      string
	PodSelfLink string
	PodQOSClass string
}

func onAdd(events chan Event) func(obj interface{}) {
	return func(obj interface{}) {
		pod, ok := obj.(*v1.Pod)
		if !ok {
			return
		}
		events <- Event{Type: NewPodEvent, PodUID: string(pod.UID), PodSelfLink: pod.SelfLink, PodQOSClass: string(pod.Status.QOSClass)}
		return
	}
}

func onUpdate(events chan Event) func(oldObj, newObj interface{}) {
	return func(oldObj, newObj interface{}) {
		// do nothing
		return
		// oldPod, ok := oldObj.(*v1.Pod)
		// if !ok {
		// 	return
		// }
		// newPod, ok := newObj.(*v1.Pod)
		// if !ok {
		// 	return
		// }
	}
}

func onDelete(events chan Event) func(obj interface{}) {
	return func(obj interface{}) {
		pod, ok := obj.(*v1.Pod)
		if !ok {
			return
		}
		events <- Event{Type: DeletePodEvent, PodUID: string(pod.UID), PodSelfLink: pod.SelfLink, PodQOSClass: string(pod.Status.QOSClass)}
		return
	}
}
