CREATE TABLE IF NOT EXISTS outbox
(
    id           SERIAL PRIMARY KEY,
    aggregate_id BIGINT    NOT NULL,
    payload      JSONB     NOT NULL,
    created_at   TIMESTAMP NOT NULL DEFAULT now(),
    processed_at TIMESTAMP NULL
);