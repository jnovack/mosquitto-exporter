PACKAGE := $(shell git remote get-url --push origin | sed -E 's/.+[@|/](.+)\.(.+).git/\1.\2/' | sed 's/\:/\//')
APPLICATION := $(shell basename `pwd`)
BUILD_RFC3339 := $(shell date -u +"%Y-%m-%dT%H:%M:%S+00:00")
COMMIT := $(shell git rev-parse HEAD)
VERSION := $(shell git describe --tags)

GO_LDFLAGS := "-w -s \
	-X github.com/jnovack/go-version.Application=${APPLICATION} \
	-X github.com/jnovack/go-version.BuildDate=${BUILD_RFC3339} \
	-X github.com/jnovack/go-version.Revision=${COMMIT} \
	-X github.com/jnovack/go-version.Version=${VERSION} \
	"

all: build

.PHONY: build
build:
	go build -o bin/${APPLICATION} -ldflags $(GO_LDFLAGS) ${PACKAGE}