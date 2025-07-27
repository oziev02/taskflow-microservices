FROM golang:1.23.6-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o event-logger ./cmd/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/event-logger .

CMD ["./event-logger"]
