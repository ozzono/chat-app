.PHONY: swag start start-docker

swag:
	# docker run --rm -v /home/hugo/Projects/snippets/jobsity/chat-app:/code -w /code ghcr.io/swaggo/swag:latest init -g cmd/main.go -o doc
	swag init -g ./cmd/main.go -o docs

start-go: swag
	go run ./cmd/

start: swag
	docker-compose up --build

test: test-controller test-repo test-queue

test-controller:
    @echo "Running tests in internal/controller..."
    @go test ./internal/controller/... -v

test-repo:
    @echo "Running tests in internal/repo..."
    @go test ./internal/repo/... -v

test-queue:
    @echo "Running tests in pkg/queue..."
    @go test ./pkg/queue/... -v