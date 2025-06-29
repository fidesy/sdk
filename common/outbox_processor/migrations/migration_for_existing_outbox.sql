-- +goose Up
-- +goose StatementBegin
ALTER TABLE table_outbox
    ADD COLUMN created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    ADD COLUMN sent_at TIMESTAMP;

CREATE INDEX table_outbox_sent_at_null_idx ON table_outbox (sent_at) WHERE sent_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS table_outbox_sent_at_null_idx;

ALTER TABLE table_outbox
    DROP COLUMN IF EXISTS created_at,
    DROP COLUMN IF EXISTS sent_at;
-- +goose StatementEnd
