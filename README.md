# cgnet

## Components

* `bpf/` - contains actual ebpf tracing code and the Go code to interact with it
* `manifests/deploy` - has Kubernetes configuration files to run `cgnet` as a Prometheus exporter
* `manifests/example` - has configuration files to set up a cluster for testing the exporter

## Dependencies

This repository experiments with [Daniel Mack](https://github.com/zonque)'s [eBPF hooks for
cgroups](https://github.com/torvalds/linux/commit/ca89fa77b4488ecf2e3f72096386e8f3a58fe2fc).

* Linux v4.10-rc
* `CONFIG_CGROUP_BPF=y`

## Vendoring

We use [dep](https://github.com/golang/dep).
