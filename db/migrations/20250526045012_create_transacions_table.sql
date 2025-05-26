-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY,
    status VARCHAR(100) NOT NULL,
    phone_number VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    is_member BOOLEAN NOT NULL DEFAULT FALSE,
    total_quantity BIGINT NOT NULL,
    total_product_amount BIGINT NOT NULL,
    total_amount BIGINT NOT NULL,
    v_payment_id VARCHAR(100) NOT NULL,
    v_payment_redirect_url VARCHAR(255) NOT NULL,
    v_transaction_id VARCHAR(100) NOT NULL,
    discount_percentage BIGINT NOT NULL,
    change_money BIGINT NOT NULL,
    payment_type VARCHAR(100) NOT NULL,
    total_money BIGINT NOT NULL,
    table_number BIGINT NOT NULL,
    total_product_capital_price BIGINT NOT NULL,
    tax_amount BIGINT NOT NULL,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_transactions_active  ON transactions(deleted_at) WHERE deleted_at IS NULL;

-- Fungsi untuk memperbarui kolom `updated_at`
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW(); -- Simpan timestamp langsung dengan tipe TIMESTAMPTZ
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger untuk tabel `transactions`
CREATE TRIGGER set_updated_at_transactions
BEFORE UPDATE ON transactions
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_transactions_active;
DROP TRIGGER IF EXISTS set_updated_at_transactions ON transactions;
DROP TABLE IF EXISTS transactions;
-- +goose StatementEnd
