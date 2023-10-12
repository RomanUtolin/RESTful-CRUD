.PHONY: build
build:
		go build -v ./cmd/app

.PHONY: test
test:
		go test -v -race -timeout 30s ./...

.PHONY: lint
lint: $(GOLANGCI) ## Runs golangci-lint with predefined configuration
	@echo "Applying linter"
	golangci-lint version
	golangci-lint run -c .golangci.yaml ./...
.DEFAULT_GOAL := build