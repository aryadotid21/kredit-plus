-- +goose Up
-- +goose StatementBegin
CREATE TABLE customer_tokens (
    id serial PRIMARY KEY,
    customer_id integer UNIQUE NOT NULL REFERENCES customers(id),
    access_token text NOT NULL,
    refresh_token text NOT NULL,
    user_agent varchar(255),
    ip_address varchar(45),
    access_token_expired_at timestamptz,
    refresh_token_expired_at timestamptz,
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);

CREATE INDEX idx_customer_tokens_customer_id ON customer_tokens (customer_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE customer_tokens;
-- +goose StatementEnd
