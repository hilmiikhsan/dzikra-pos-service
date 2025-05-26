-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS transaction_items (
    id SERIAL PRIMARY KEY,
    quantity BIGINT NOT NULL,
    total_amount BIGINT NOT NULL,
    product_name VARCHAR(100) NOT NULL,
    product_price BIGINT NOT NULL,
    transaction_id UUID NOT NULL,
    product_id INT NOT NULL,
    product_capital_price BIGINT NOT NULL,
    total_amount_capital_price BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

ALTER TABLE transaction_items ADD CONSTRAINT fk_transaction_items_transactions FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE transaction_items ADD CONSTRAINT fk_transaction_items_products FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT ON UPDATE CASCADE;

CREATE INDEX idx_transaction_items_active  ON transaction_items(deleted_at) WHERE deleted_at IS NULL;

-- Fungsi untuk memperbarui kolom `updated_at`
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW(); -- Simpan timestamp langsung dengan tipe TIMESTAMPTZ
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger untuk tabel `transaction_items`
CREATE TRIGGER set_updated_at_transaction_items
BEFORE UPDATE ON transaction_items
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_transaction_items_active;
DROP TRIGGER IF EXISTS set_updated_at_transaction_items ON transaction_items;
DROP TABLE IF EXISTS transaction_items;
-- +goose StatementEnd
