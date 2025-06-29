package outbox_processor

import "time"

type Message struct {
	ID        int64      `db:"id"`
	Message   string     `db:"message"`
	CreatedAt time.Time  `db:"created_at"`
	SentAt    *time.Time `db:"sent_at"`
}
