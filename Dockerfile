FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0

RUN go build -o bin/service ./cmd/main.go

# ---- Runtime ----
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/bin/service /app/bin/service
COPY --from=builder /app/.env /app/.env

WORKDIR /app/bin

CMD ["./service"]
