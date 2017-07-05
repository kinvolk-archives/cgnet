package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	port      int32  = 9101
	namespace string = "cgnet_pod"
)

type PodMetrics struct {
	TotalNumberPods prometheus.Gauge
	IncomingPackets *prometheus.CounterVec
	OutgoingPackets *prometheus.CounterVec
	// ...
}

var podmetrics = PodMetrics{
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

func init() {
	prometheus.MustRegister(podmetrics.TotalNumberPods)
	prometheus.MustRegister(podmetrics.IncomingPackets)
	prometheus.MustRegister(podmetrics.OutgoingPackets)
}

func serveMetrics() {
	addr := fmt.Sprintf(":%d", port)
	log.Printf("started serving metrics on %s", addr)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(addr, nil))
}
