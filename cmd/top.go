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
	"time"

	"github.com/kinvolk/cgnet/bpf"
	"github.com/spf13/cobra"
)

var topCmd = &cobra.Command{
	Use:   "top [CGROUP...]",
	Short: "Show network stats per control group",
	Run:   cmdTop,
}

func cmdTop(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		os.Exit(0)
	}

	cgPath := args[0]

	ctl, err := bpf.Attach(cgPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to load bpf:", err)
		os.Exit(1)
	}

	var packets, bytes uint64
	ctl.SetPacketsHandler(func(v uint64) error {
		packets = v
		return nil
	})
	ctl.SetBytesHandler(func(v uint64) error {
		bytes = v
		return nil
	})

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	go ctl.Run(ctx)
	defer ctl.Stop()

	term := make(chan os.Signal)
	signal.Notify(term, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-term:
			return
		case <-time.After(1 * time.Second):
			fmt.Fprintf(os.Stdout, "\rcgroup received %d packets (%d bytes)", packets, bytes)
		}
	}
}

func init() {
	RootCmd.AddCommand(topCmd)
}
