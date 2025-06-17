CREATE TABLE IF NOT EXISTS payment_outbox
(
    id           SERIAL PRIMARY KEY,
    message_key  TEXT UNIQUE NOT NULL,
    event_type   TEXT        NOT NULL,
    message      TEXT,
    payload      JSONB       NOT NULL,
    created_at   TIMESTAMP   NOT NULL DEFAULT now(),
    processed_at TIMESTAMP   NULL
);