.PHONY: all build
pwd := $(shell pwd)

all: run

run:
	GOPATH=$(pwd)/_vendor go run collectd_sidekiq.go $(ARGS)

build:
	GOOS=linux GOARCH=amd64 GOPATH=$(pwd)/_vendor go build -o collectd_sidekiq.linux_amd64 collectd_sidekiq.go
