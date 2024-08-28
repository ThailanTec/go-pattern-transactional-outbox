package domain

import "time"

type Order struct {
	ID        int       `db:"id"`
	Item      string    `db:"item"`
	Quantity  int       `db:"quantity"`
	CreatedAt time.Time `db:"created_at"`
}

type OutboxMessage struct {
	ID        int       `db:"id"`
	Payload   string    `db:"payload"`
	Processed bool      `db:"processed"`
	CreatedAt time.Time `db:"created_at"`
}
