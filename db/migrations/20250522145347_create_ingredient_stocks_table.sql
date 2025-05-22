-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS ingredient_stocks (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    unit VARCHAR(100) NOT NULL,
    price_per_amount_stock BIGINT NOT NULL,
    required_stock BIGINT NOT NULL,
    amount_price_required_stock BIGINT NOT NULL,
    amount_stock_per_price NUMERIC(10,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_ingredient_stocks_active  ON ingredient_stocks(deleted_at) WHERE deleted_at IS NULL;

-- Fungsi untuk memperbarui kolom `updated_at`
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW(); -- Simpan timestamp langsung dengan tipe TIMESTAMPTZ
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger untuk tabel `ingredient_stocks`
CREATE TRIGGER set_updated_at_ingredient_stocks
BEFORE UPDATE ON ingredient_stocks
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_ingredient_stocks_active;
DROP TRIGGER IF EXISTS set_updated_at_ingredient_stocks ON ingredient_stocks;
DROP TABLE IF EXISTS ingredient_stocks;
-- +goose StatementEnd
