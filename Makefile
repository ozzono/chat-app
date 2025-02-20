
.PHONY: swag start

swag:
	swag init -g ./cmd/main.go -o docs

start: swag
	go run ./cmd/
