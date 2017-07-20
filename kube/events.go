package kube

import (
	log "github.com/inconshreveable/log15"
	"k8s.io/client-go/pkg/api/v1"
)

type Event int

const (
	NewPodEvent Event = iota
	UpdatePodEvent
	DeletePodEvent
)

func onAdd(events chan Event) func(obj interface{}) {
	funclog := log.New("func", "onAdd")
	return func(obj interface{}) {
		pod, ok := obj.(*v1.Pod)
		if !ok {
			funclog.Error("unable to assert type")
			return
		}
		events <- NewPodEvent
		funclog.Info(pod.ObjectMeta.SelfLink)
	}
}

func onUpdate(events chan Event) func(oldObj, newObj interface{}) {
	// funclog := log.New("func", "onUpdate")
	return func(oldObj, newObj interface{}) {
		// do nothing
		return
		// oldPod, ok := oldObj.(*v1.Pod)
		// if !ok {
		// 	funclog.Error("unable to assert type")
		// 	return
		// }
		// newPod, ok := newObj.(*v1.Pod)
		// if !ok {
		// 	funclog.Error("unable to assert type")
		// 	return
		// }
		// funclog.Info(fmt.Sprintf("old: %s, new: %s", oldPod.ObjectMeta.SelfLink, newPod.ObjectMeta.SelfLink))
	}
}

func onDelete(events chan Event) func(obj interface{}) {
	funclog := log.New("func", "onDelete")
	return func(obj interface{}) {
		pod, ok := obj.(*v1.Pod)
		if !ok {
			funclog.Error("unable to assert type")
			return
		}
		events <- DeletePodEvent
		funclog.Info(pod.ObjectMeta.SelfLink)
	}
}
