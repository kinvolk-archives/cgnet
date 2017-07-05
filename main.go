package main

func main() {
	stop := make(chan struct{})
	defer close(stop)

	go RunPodInformer(stop)

	// TODO 
	// * install bpf program on every 'new pod' event
	// * query the bpf maps to retrieve data
	// * update podmetrics with data

	serveMetrics()
}
