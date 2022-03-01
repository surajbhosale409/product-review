SHELL = /bin/sh

.PHONY: vendor fmt lint vet test clean build

tools:
	sh -c "$$(wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh || echo exit 2)" -- -b $(shell go env GOPATH)/bin $(GOLANGCI_LINT_VERSION)
	go build -o build/seed-db tools/seed-db.go

lint:
	golangci-lint run ./...

vet:
	go vet ./...

fmt:
	@echo $(shell go fmt ./...)

test:
	go test ./...

clean:
	rm -rf build

build: clean
	mkdir -p build
	go build -o build/product-review cmd/main/product-review.go

install:
	go install cmd/main/product-review.go

vendor:
	go mod tidy
	go mod vendor

image:
	docker build -t product-review:latest .

up:
	docker-compose up -d

down:
	docker-compose down

seed-db: tools up
	build/seed-db

run: build up
	build/product-review serve