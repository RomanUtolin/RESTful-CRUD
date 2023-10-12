.PHONY: build
build:
		go build -v ./cmd/app

.PHONY: test
test:
		go test -v -race -timeout 30s ./...

.PHONY: start
start:
	docker build -t app-server-crud .
	docker compose up

.DEFAULT_GOAL := build