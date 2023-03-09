VERSION=$(shell git describe --tags --long --dirty 2>/dev/null)

## we must have tagged the repo at least once for VERSION to work
ifeq ($(VERSION),)
	VERSION = UNKNOWN
endif

build:
	GOARCH=amd64 GOOS=linux go build -ldflags "-X main.version=${VERSION}" -o gocdcl ./cmd

.PHONY: lint
lint:
	golangci-lint run

.PHONY: committed
committed:
	@git diff --exit-code > /dev/null || (echo "** COMMIT YOUR CHANGES FIRST **"; exit 1)