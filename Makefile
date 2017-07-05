CONTAINER=10.0.0.240/cgnet-exporter
BIN=cgnet-exporter
GOOS=linux

.PHONY: all clean build container

all: $(BIN)

build: $(BIN)
$(BIN):
	GOOS=$(GOOS) go build -o $@ .

container: $(BIN)
	docker build -t $(CONTAINER):latest .
	docker push $(CONTAINER):latest

clean:
	rm -rf $(BIN)

dist-clean:
	docker rmi $(CONTAINER):latest
