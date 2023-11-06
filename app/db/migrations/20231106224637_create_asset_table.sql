-- +goose Up
-- +goose StatementBegin
CREATE TYPE "enum_assets_type" AS ENUM (
    'White Goods',
    'Motor',
    'Mobil'
);

CREATE TABLE assets (
    id serial PRIMARY KEY,
    name varchar(255) NOT NULL,
    type enum_assets_type NOT NULL,
    description text,
    price numeric(15, 2),
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE assets;
-- +goose StatementEnd
