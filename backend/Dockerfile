FROM golang:latest AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o main .

RUN go build -o parse_fias ./scripts

FROM debian:unstable-slim

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/parse_fias .
COPY --from=builder /app/logs .
COPY --from=builder /app/upload .
COPY --from=builder /app/settings.json .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/scripts/import ./scripts/import

EXPOSE 8080

CMD ["./main"]