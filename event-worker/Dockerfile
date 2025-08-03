FROM golang:1.23.6

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o event-worker ./cmd/main.go

CMD ["./event-worker"]
