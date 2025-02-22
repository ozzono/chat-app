.PHONY: swag start start-go test test-controller test-repo test-queue

swag:
	@echo "Generating Swagger documentation..."
	swag init -g ./cmd/main.go -o docs

start-go: swag
	go run ./cmd/

start: swag
	docker-compose up --build

test: test-controller test-repo test-queue test-bot

test-controller:
	@echo "Running tests in internal/controller..."
	@go test ./internal/controller/... -v

test-repo:
	@echo "Running tests in internal/repo..."
	@go test ./internal/repo/... -v

test-queue:
	@echo "Running tests in pkg/queue..."
	@go test ./pkg/queue/... -v

test-bot:
	@echo "Running tests in pkg/bot..."
	@go test ./pkg/bot/... -v