CREATE TABLE IF NOT EXISTS accounts
(
    id         serial PRIMARY KEY,
    user_id    BIGINT         NOT NULL UNIQUE,
    balance    NUMERIC(18, 2) NOT NULL DEFAULT 0,
    updated_at TIMESTAMP      NOT NULL DEFAULT now()

);