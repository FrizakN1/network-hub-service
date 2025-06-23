FROM golang:latest AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o network-hub-service .

FROM debian:unstable-slim

WORKDIR /app

COPY --from=builder /app/network-hub-service .
COPY --from=builder /app/logs ./logs
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/upload ./upload

RUN touch .env

EXPOSE 8080

CMD ["./network-hub-service"]