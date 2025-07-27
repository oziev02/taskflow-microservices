FROM golang:1.23.6-alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o task-service ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/task-service .

COPY .env .env

EXPOSE 8080

CMD ["./task-service"]
