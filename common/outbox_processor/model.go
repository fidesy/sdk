package outbox_processor

type Message struct {
	ID      int64  `db:"id"`
	Message string `db:"message"`
}
