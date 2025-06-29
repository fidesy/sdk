-- +goose Up
-- +goose StatementBegin
CREATE TABLE table_outbox (
                                 id SERIAL PRIMARY KEY,
                                 message TEXT,
                                 created_at TIMESTAMP DEFAULT NOW() NOT NULL,
                                 sent_at TIMESTAMP
);

CREATE INDEX table_outbox_sent_at_null_idx ON table_outbox (sent_at) WHERE sent_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE table_outbox;
-- +goose StatementEnd
