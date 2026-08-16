.PHONY: build test lint run clean

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

build:
	go build $(LDFLAGS) -o bin/chain-health ./cmd/health

test:
	go test -race -cover ./...

lint:
	golangci-lint run ./...

run: build
	./bin/chain-health --config config.yaml

docker:
	docker build -t chain-health:$(VERSION) .

clean:
	rm -rf bin/
