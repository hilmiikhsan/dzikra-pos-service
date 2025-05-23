-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tax (
    id SERIAL PRIMARY KEY,
    tax BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_tax_active  ON tax(deleted_at) WHERE deleted_at IS NULL;

-- Fungsi untuk memperbarui kolom `updated_at`
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW(); -- Simpan timestamp langsung dengan tipe TIMESTAMPTZ
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger untuk tabel `tax`
CREATE TRIGGER set_updated_at_tax
BEFORE UPDATE ON tax
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_tax_active;
DROP TRIGGER IF EXISTS set_updated_at_tax ON tax;
DROP TABLE IF EXISTS tax;
-- +goose StatementEnd
