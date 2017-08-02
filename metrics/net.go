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

package metrics

import "github.com/prometheus/client_golang/prometheus"

type PodMetrics struct {
	TotalNumberPods prometheus.Gauge
	IncomingPackets *prometheus.CounterVec
	OutgoingPackets *prometheus.CounterVec
	// ...
}

var globalPodMetrics = PodMetrics{
	TotalNumberPods: prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name:      "total_number_pods",
			Namespace: namespace,
			Help:      "Total number of pods in the cluster",
		},
	),
	IncomingPackets: prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "packets_incoming_total",
			Namespace: namespace,
			Help:      "Total number of incoming packets.",
		},
		[]string{"pod_name"},
	),
	OutgoingPackets: prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "packets_outgoing_total",
			Namespace: namespace,
			Help:      "Total number of outgoing packets.",
		},
		[]string{"pod_name"},
	),
}

func TotalNum() prometheus.Gauge {
	return globalPodMetrics.TotalNumberPods
}

func SetOutgoingPackets(pod string, packets float64) {
	globalPodMetrics.OutgoingPackets.With(prometheus.Labels{"pod_name": pod}).Add(packets)
}

func SetIncomingPackets(pod string, packets float64) {
	globalPodMetrics.IncomingPackets.With(prometheus.Labels{"pod_name": pod}).Add(packets)
}
