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
)

var (
	version      string
	printVersion bool
)

var RootCmd = &cobra.Command{
	Use:   "cgnet",
	Short: "cgroup network statistics",
	Run: func(cmd *cobra.Command, _ []string) {
		if printVersion {
			fmt.Fprintln(os.Stdout, cmd.Name(), version)
			os.Exit(0)
		}
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&printVersion, "version", "V", false, "print version information")
}
