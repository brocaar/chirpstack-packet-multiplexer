.PHONY: build test dist snapshot dev-requirements requirements
VERSION := $(shell git describe --always |sed -e "s/^v//")

build:
	mkdir -p build
	go build -ldflags "-s -w -X main.version=$(VERSION)" -o build/lora-packet-multiplexer cmd/lora-packet-multiplexer/main.go

clean:
	rm -rf build
	rm -rf dist
	rm -rf docs/public

test:
	go vet ./...
	go test -v ./...

dist:
	goreleaser
	mkdir -p dist/upload/tar
	mkdir -p dist/upload/deb
	mv dist/*.tar.gz dist/upload/tar
	mv dist/*.deb dist/upload/deb

snapshot:
	goreleaser --snapshot

dev-requirements:
	go get -u golang.org/x/tools/cmd/stringer
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/goreleaser/goreleaser
	go get -u github.com/goreleaser/nfpm

requirements:
	dep ensure -v
