VERSION=$(shell git describe --tags)

build:
	go build -ldflags "-X main.version=$(VERSION)"

install:
	go install -ldflags "-X main.version=$(VERSION)"
