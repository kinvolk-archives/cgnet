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

package kube

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const DefaultCgroupRoot string = "/sys/fs/cgroup/systemd/kubepods"

func BuildConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}

// We are using the default root for now.
// It looks like querying of the kubelet config is part of
// this PR for dynamic kubelet configuration:
//
// https://github.com/kubernetes/features/issues/281
// https://github.com/kubernetes/kubernetes/pull/46254
//
// TODO: investigate ^ or
// https://github.com/kubernetes/community/blob/master/contributors/design-proposals/dynamic-kubelet-configuration.md#monitoring-configuration-status
// there is some mention of the `configz` endpoint
func GetCgroupRoot(_ *rest.Config) (string, error) {
	return DefaultCgroupRoot, nil
}
