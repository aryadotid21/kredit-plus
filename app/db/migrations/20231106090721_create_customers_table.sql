-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE customers (
    id serial PRIMARY KEY,
    uuid uuid UNIQUE DEFAULT uuid_generate_v4(),
    email varchar(255) UNIQUE CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    phone varchar(20) UNIQUE CHECK (phone ~* '^[0-9]+$'),
    password text,
    last_login timestamptz,
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);

CREATE INDEX idx_customers_uuid ON customers (uuid);
CREATE INDEX idx_customers_email ON customers (email);
CREATE INDEX idx_customers_phone ON customers (phone);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE customers;

DROP EXTENSION IF EXISTS "uuid-ossp";
-- +goose StatementEnd
