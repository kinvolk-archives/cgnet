package kube

import (
	"log"

	"k8s.io/api/core/v1"
)

type Event int

const (
	NewPodEvent Event = iota
	DeletePodEvent
)

func emitEvent(eChan chan Event, e Event) func(obj interface{}) {
	return func(obj interface{}) {
		_, ok := obj.(*v1.Pod)
		if !ok {
			log.Printf("unexpected object type: %#v\n", obj)
		}
		eChan <- e
	}
}
