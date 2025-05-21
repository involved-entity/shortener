FROM golang:1.24.2-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN GOPROXY=direct go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o shortener ./cmd/shortener/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/shortener .
COPY --from=builder /app/configs ./configs

EXPOSE 8000

ENV CONFIG_PATH=config/local.docker.yml

CMD ["./shortener"]
