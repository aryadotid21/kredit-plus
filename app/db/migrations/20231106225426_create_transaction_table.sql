-- +goose Up
-- +goose StatementBegin
CREATE TABLE transactions (
    id serial PRIMARY KEY,
    uuid uuid DEFAULT uuid_generate_v4(),
    customer_id integer NOT NULL REFERENCES customers(id),
    asset_id integer NOT NULL REFERENCES assets(id),
    contract_number varchar(255) NOT NULL,
    otr_amount numeric(15, 2) NOT NULL,
    admin_fee numeric(15, 2) NOT NULL,
    installment_amount numeric(15, 2) NOT NULL,
    interest_amount numeric(15, 2) NOT NULL,
    sales_channel varchar(255) NOT NULL,
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);

CREATE INDEX idx_transactions_customer_id ON transactions (customer_id);
CREATE INDEX idx_transactions_asset_id ON transactions (asset_id);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE transactions;
-- +goose StatementEnd
