.PHONY: swag start start-go test test-go-controller test-go-repo test-go-queue

swag:
	@echo "Generating Swagger documentation..."
	swag init -g ./cmd/main.go -o docs

start-go: swag
	go run ./cmd/

start:
	docker-compose up -d api

test:
	docker-compose up -d test

test-go: test-go-controller test-go-repo test-go-queue test-go-bot

test-go-controller:
	@echo "Running tests in internal/controller..."
	@go test ./internal/controller/... -v

test-go-repo:
	@echo "Running tests in internal/repo..."
	@go test ./internal/repo/... -v

test-go-queue:
	@echo "Running tests in pkg/queue..."
	@go test ./pkg/queue/... -v

test-go-bot:
	@echo "Running tests in pkg/bot..."
	@go test ./pkg/bot/... -v