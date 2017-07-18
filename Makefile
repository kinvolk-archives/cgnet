CONTAINER=10.0.0.240/cgnet-exporter
BIN=cgnet-exporter
GOOS=linux

.PHONY: all clean build container manifest

all: $(BIN)

build: $(BIN)
$(BIN):
	GOOS=$(GOOS) go build -o $@ .

container: $(BIN)
	docker build -t $(CONTAINER):latest .
	docker push $(CONTAINER):latest

manifest: deploy/all-in-one.yaml
deploy/all-in-one.yaml:
	@make -C deploy/

clean:
	rm -rf $(BIN)
	@make -C deploy/ clean

dist-clean:
	docker rmi $(CONTAINER):latest
