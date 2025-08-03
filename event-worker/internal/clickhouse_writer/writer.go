package clickhouse_writer

import (
	"context"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type Writer struct {
	conn clickhouse.Conn
}

func NewClickHouseWriter(addr string) *Writer {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{addr},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
	})
	if err != nil {
		log.Fatalf("failed to connect to ClickHouse: %v", err)
	}
	return &Writer{conn: conn}
}

func (w *Writer) HandleEvent(topic string, message []byte) error {
	ctx := context.Background()
	query := `
		INSERT INTO task_events (timestamp, topic, payload)
		VALUES (?, ?, ?)
	`
	err := w.conn.Exec(ctx, query, time.Now(), topic, string(message))
	if err != nil {
		log.Printf("failed to write event to ClickHouse: %v", err)
	}
	return err
}
