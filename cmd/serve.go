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

package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/inconshreveable/log15"
	"github.com/spf13/cobra"

	"github.com/kinvolk/cgnet/kube"
	"github.com/kinvolk/cgnet/metrics"
)

var (
	metricsPort int
	kubeconfig  string
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve metrics to Prometheus",
	Run:   cmdServe,
}

func cmdServe(cmd *cobra.Command, args []string) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	cfg, err := kube.BuildConfig(kubeconfig)
	if err != nil {
		log.Error("error building config", "err", err)
		return
	}

	events := make(chan kube.Event)
	go kube.WatchPodEvents(ctx, cancelFunc, cfg, events)

	addr := fmt.Sprintf(":%d", metricsPort)
	go metrics.Serve(ctx, addr)

	// TODO
	// * install bpf program on every 'new pod' event
	// * query the bpf maps to retrieve data
	// * update podmetrics with data

	term := make(chan os.Signal)
	signal.Notify(term, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-term:
			return
		case <-ctx.Done():
			return
		case e := <-events:
			switch e {
			case kube.NewPodEvent:
				metrics.GlobalPodMetrics.TotalNumberPods.Add(1)
			case kube.DeletePodEvent:
				metrics.GlobalPodMetrics.TotalNumberPods.Sub(1)
			}
		}
	}
}

func init() {
	RootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntVarP(&metricsPort, "port", "p", 9101, "metrics port")
	serveCmd.Flags().StringVarP(&kubeconfig, "kubeconfig", "k", "", "path to kubeconfig file. Only required if out-of-cluster.")
}
