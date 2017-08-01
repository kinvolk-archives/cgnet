# cgnet

`cgnet` uses eBPF to gather network statistics from cgroups.

```
$ make

# to test the cli tool
$ ./cgnet top <path to cgroup>

# to run the exporter outside of k8s
$ ./cgnet export --kubeconfig <path to kubeconfig>

# to run the exporter as a DaemonSet in the cluster
$ kubectl apply -f manifests/deploy/all-in-one.yaml
```

## cgnet top

Uses the BPF program to monitor the cgroup specified on the command line.

## cgnet export

Interacts with the Kubernetes API to retrieve a list of running pods and installs the BPF program for their cgroups.
It then exports the data for Prometheus.

We provide the following to deploy and test the exporter:

* `Dockerfile` - builds a simple minimal container running `cgnet export`
* `manifests/deploy` - has Kubernetes manifests to run the container with the exporter
* `manifests/example` - has configuration files to set up a cluster for testing the exporter

## Dependencies

This repository experiments with [Daniel Mack](https://github.com/zonque)'s [eBPF hooks for
cgroups](https://github.com/torvalds/linux/commit/ca89fa77b4488ecf2e3f72096386e8f3a58fe2fc).

* Linux v4.10-rc
* `CONFIG_CGROUP_BPF=y`

## Vendoring

We use [dep](https://github.com/golang/dep).
