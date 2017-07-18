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
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/kinvolk/cgnet/bpf"
)

var RootCmd = &cobra.Command{
	Use:   "cgnet <cgroup-path>",
	Short: "cgroup-ebpf experiment",
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] == "" {
			return cmd.Usage()
		}

		if err := bpf.Setup(args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to load bpf: %s", err)
			os.Exit(1)
		}

		quit := make(chan struct{})
		defer close(quit)

		if err := bpf.UpdateLoop(quit); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to load bpf: %s", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
