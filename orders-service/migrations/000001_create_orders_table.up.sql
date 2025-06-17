CREATE TABLE orders
(
    id         serial PRIMARY KEY,
    user_id    serial NOT NULL,
    status     text   not null DEFAULT 'New',
    message    text            DEFAULT '',
    amount     TEXT   NOT NULL,
    created_at TIMESTAMP       DEFAULT now()
);
