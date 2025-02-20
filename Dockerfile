# Stage 1: Build the binary
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . ./

# Generate Swagger docs
RUN swag init -g ./cmd/main.go -o docs

RUN go build -o /chat-app ./cmd/main.go

# Stage 2: Run the binary
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /chat-app .
COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./chat-app"]