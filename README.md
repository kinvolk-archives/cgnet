# cgnet

## Components

* `ds.yaml` - runs cgnet as a Kubernetes DaemonSet
* `cgt-agent` - collect Pod data from Kubernetes API, install eBPF programs and expose data
* `bpf/` - contains actual ebpf tracing code
* `prometheus.yaml` - example configuration for prometheus

### What is does:

- [ ] Detect pods on the node
- [ ] Subscribe to k8s new pod event
- [ ] Run `bpf-tracer` on all nodes
- [ ] Read bpf-map
- [ ] Expose data to prometheus

## Dependencies

This repository experiments with [Daniel Mack](https://github.com/zonque)'s [eBPF hooks for
cgroups](https://github.com/torvalds/linux/commit/ca89fa77b4488ecf2e3f72096386e8f3a58fe2fc).

* Linux v4.10-rc
* `CONFIG_CGROUP_BPF=y`

## Vendoring

We use [gvt](https://github.com/FiloSottile/gvt).
