FROM scratch

COPY cgnet-exporter /cgnet-exporter

EXPOSE 9101
WORKDIR /data
ENTRYPOINT ["/cgnet-exporter"]
