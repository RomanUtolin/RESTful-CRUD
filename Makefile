.PHONY: build
build:
		go build -v ./cmd/app

.PHONY: test
test:
		go test -v -timeout 30s ./...

.PHONY: start
start:
	docker build -t app-server-crud .
	docker compose up

.PHONY: testDb
testDb:
	docker compose --file compose-testDb.yaml up

.PHONY: mockery
mockery:
	go run github.com/vektra/mockery/v2@v2.36.0 --all

.DEFAULT_GOAL := build