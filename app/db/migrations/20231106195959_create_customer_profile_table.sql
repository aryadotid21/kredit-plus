-- +goose Up
-- +goose StatementBegin

CREATE TABLE customer_profiles (
    id serial PRIMARY KEY,
    customer_id integer UNIQUE NOT NULL REFERENCES customers(id),
    NIK varchar(20),
    full_name varchar(255),
    legal_name varchar(255),
    place_of_birth varchar(255),
    date_of_birth date,
    salary numeric(15, 2),
    ktp_image text,
    selfie_image text,
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);

CREATE INDEX idx_customer_id ON customer_profiles (customer_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE customer_profiles;

-- +goose StatementEnd
