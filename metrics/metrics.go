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

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const namespace string = "cgnet_pod"

func init() {
	prometheus.MustRegister(globalPodMetrics.TotalNumberPods)
	prometheus.MustRegister(globalPodMetrics.IncomingPackets)
	prometheus.MustRegister(globalPodMetrics.OutgoingPackets)
}

func Serve(ctx context.Context, addr string) {
	http.Handle("/metrics", promhttp.Handler())
	srv := http.Server{
		Addr:    addr,
		Handler: http.DefaultServeMux,
	}
	go srv.ListenAndServe()

	<-ctx.Done()

	toCtx, cancelFunc := context.WithTimeout(ctx, 2*time.Second)
	defer cancelFunc()

	if err := srv.Shutdown(toCtx); err != nil {
		panic(err)
	}
}
