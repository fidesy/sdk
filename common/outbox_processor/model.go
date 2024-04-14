package outbox_processor

type OutboxMessage struct {
	ID      int64  `db:"id"`
	Message string `db:"message"`
}
