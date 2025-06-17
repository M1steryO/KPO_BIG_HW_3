CREATE TABLE IF NOT EXISTS inbox
(
    id           SERIAL PRIMARY KEY,
    message_key  TEXT UNIQUE NOT NULL,
    payload      JSONB       NOT NULL,
    created_at   TIMESTAMP   NOT NULL DEFAULT now(),
    processed_at TIMESTAMP   NULL
);