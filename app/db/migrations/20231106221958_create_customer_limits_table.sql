-- +goose Up
-- +goose StatementBegin
CREATE TABLE customer_limits (
    id serial PRIMARY KEY,
    customer_id integer NOT NULL REFERENCES customers(id),
    tenor integer,
    limit_amount numeric(15, 2),
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_customer_limits_customer_id_tenor ON customer_limits (customer_id, tenor);
CREATE INDEX idx_customer_limits_customer_id ON customer_limits (customer_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE customer_limits;
-- +goose StatementEnd