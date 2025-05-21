FROM golang:1.24.2-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o worker ./cmd/worker/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/worker .
COPY --from=builder /app/configs ./configs
CMD ["./worker"]
