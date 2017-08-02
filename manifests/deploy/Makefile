TARGET:=all-in-one.yaml

.PHONY: all clean

all: $(TARGET)
$(TARGET):
	cat cgnet-exporter-ds.yaml > $@
	@echo "---" >> $@
	cat cgnet-exporter-svc.yaml >> $@
	@echo "---" >> $@
	cat cgnet-exporter-svc-monitor.yaml >> $@

clean: $(TARGET)
	rm -rf $<
