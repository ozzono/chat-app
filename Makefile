.PHONY: swag start start-docker

swag:
	# docker run --rm -v /home/hugo/Projects/snippets/jobsity/chat-app:/code -w /code ghcr.io/swaggo/swag:latest init -g cmd/main.go -o doc
	swag init -g ./cmd/main.go -o docs

start-go: swag
	go run ./cmd/

start: swag
	docker-compose up --build
