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
	"path"
	"strings"
	"syscall"

	log "github.com/inconshreveable/log15"
	"github.com/spf13/cobra"

	"github.com/kinvolk/cgnet/bpf"
	"github.com/kinvolk/cgnet/kube"
	"github.com/kinvolk/cgnet/metrics"
)

var (
	metricsPort      int
	kubeconfig       string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Serve metrics to Prometheus",
	Run:   cmdExport,
}

func cmdExport(cmd *cobra.Command, args []string) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// needs to be run before any API interaction
	cfg, err := kube.BuildConfig(kubeconfig)
	if err != nil {
		log.Error("error building config", "err", err)
		return
	}

	cgroupRoot, err := kube.GetCgroupRoot(cfg)
	if err != nil {
		log.Error("error retrieving cgroup root for cluster", "err", err)
		return
	}
	log.Debug("cgroup root", "path", cgroupRoot)

	events := make(chan kube.Event)
	go kube.WatchPodEvents(ctx, cancelFunc, cfg, events)

	addr := fmt.Sprintf(":%d", metricsPort)
	go metrics.Serve(ctx, addr)

	term := make(chan os.Signal)
	signal.Notify(term, syscall.SIGINT, syscall.SIGTERM)

	var ctlMap map[string]*bpf.Controller
	for {
		select {
		case <-term:
			return
		case <-ctx.Done():
			return
		case e := <-events:
			switch e.Type {
			case kube.NewPodEvent:
				metrics.TotalNum().Add(1)
				log.Debug("new pod", "pod", e.PodSelfLink, "cgroup", buildCgroupPath(cgroupRoot, e.PodUID, e.PodQOSClass))
				// FIXME: https://github.com/kinvolk/cgnet/issues/16
				cgPath := buildCgroupPath(cgroupRoot, e.PodUID, e.PodQOSClass)
				bpfController, err := bpf.Attach(cgPath)
				if err != nil {
					continue
				}

				bpfController.SetPacketsHandler(func(v uint64) error {
					metrics.SetOutgoingPackets(path.Base(e.PodSelfLink), float64(v))
					return nil
				})
				go bpfController.Run(ctx)
				ctlMap[e.PodUID] = bpfController

			case kube.DeletePodEvent:
				metrics.TotalNum().Sub(1)
				log.Debug("pod gone", "pod", e.PodSelfLink, "cgroup", buildCgroupPath(cgroupRoot, e.PodUID, e.PodQOSClass))
				metrics.TotalNum().Sub(1)
				ctlMap[e.PodUID].Stop()
			}
		}
	}
}

func buildCgroupPath(root, uid, qosclass string) string {
	return fmt.Sprintf("%s/%s/pod%s", root, strings.ToLower(qosclass), uid)
}

func init() {
	RootCmd.AddCommand(exportCmd)

	exportCmd.Flags().IntVarP(&metricsPort, "port", "p", 9101, "metrics port")
	exportCmd.Flags().StringVarP(&kubeconfig, "kubeconfig", "k", "", "path to kubeconfig file. Only required if out-of-cluster.")
}
