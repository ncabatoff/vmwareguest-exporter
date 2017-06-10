# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.8

RUN apt-get update
RUN apt-get install -y open-vm-tools open-vm-tools-dev

#COPY . /go/src/github.com/ncabatoff/vmwareguest-exporter
#RUN go get github.com/ncabatoff/go-vmguestlib/vmguestlib
#RUN go get github.com/prometheus/client_golang/prometheus
#RUN go build github.com/ncabatoff/vmwareguest-exporter
RUN go get github.com/ncabatoff/vmwareguest-exporter
RUN go install github.com/ncabatoff/vmwareguest-exporter

USER root

ENTRYPOINT /go/bin/vmwareguest-exporter

EXPOSE 9263
