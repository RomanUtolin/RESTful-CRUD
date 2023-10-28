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

.PHONY: mockery
mockery:
	go run github.com/vektra/mockery/v2@v2.36.0 --all

.DEFAULT_GOAL := build