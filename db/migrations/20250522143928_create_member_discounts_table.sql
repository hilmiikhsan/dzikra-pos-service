-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS member_discounts (
    id SERIAL PRIMARY KEY,
    discount BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_member_discounts_active  ON member_discounts(deleted_at) WHERE deleted_at IS NULL;

-- Fungsi untuk memperbarui kolom `updated_at`
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW(); -- Simpan timestamp langsung dengan tipe TIMESTAMPTZ
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger untuk tabel `member_discounts`
CREATE TRIGGER set_updated_at_member_discounts
BEFORE UPDATE ON member_discounts
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_member_discounts_active;
DROP TRIGGER IF EXISTS set_updated_at_member_discounts ON member_discounts;
DROP TABLE IF EXISTS member_discounts;
-- +goose StatementEnd
