FROM golang AS builder
RUN apt-get update && apt-get install -y open-vm-tools open-vm-tools-dev
COPY . /go/src/github.com/ncabatoff/vmwareguest-exporter
RUN go install github.com/ncabatoff/vmwareguest-exporter

FROM debian:stretch-slim
COPY --from=builder /go/bin/vmwareguest-exporter /usr/local/bin
RUN apt-get update && apt-get install -y open-vm-tools && rm -rf /var/lib/apt/lists/*
RUN groupadd -r vmwareguest-exporter --gid=999 && useradd --no-log-init -r -g vmwareguest-exporter --uid=999 vmwareguest-exporter
USER 999
ENTRYPOINT ["/usr/local/bin/vmwareguest-exporter"]
