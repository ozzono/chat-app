# unfortunately Gos does not support CGO enabled builds
FROM debian:12-slim

RUN apt-get update && apt-get install -y \
    wget \
    gcc \
    libc6-dev \
    && rm -rf /var/lib/apt/lists/*

RUN wget -q https://go.dev/dl/go1.24.0.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz \
    && rm go1.24.0.linux-amd64.tar.gz

ENV CGO_ENABLED=1
ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /app

COPY . .

RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@latest
ENV PATH="/root/go/bin:${PATH}"
RUN swag init -g ./cmd/main.go -o docs

RUN go build -o chat-app ./cmd/main.go

EXPOSE 8080

CMD ["./chat-app"]