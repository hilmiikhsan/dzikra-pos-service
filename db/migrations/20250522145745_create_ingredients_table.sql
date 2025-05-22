-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS ingredients (
    id SERIAL PRIMARY KEY,
    unit VARCHAR(100) NOT NULL,
    cost BIGINT NOT NULL,
    recipe_id INT NOT NULL,
    ingredient_stock_id INT NOT NULL,
    required_stock NUMERIC(10,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

ALTER TABLE ingredients ADD CONSTRAINT fk_ingredients_recipe_id FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE ingredients ADD CONSTRAINT fk_ingredients_ingredient_stock_id FOREIGN KEY (ingredient_stock_id) REFERENCES ingredient_stocks(id) ON DELETE RESTRICT ON UPDATE CASCADE;

CREATE INDEX idx_ingredients_active  ON ingredients(deleted_at) WHERE deleted_at IS NULL;

-- Fungsi untuk memperbarui kolom `updated_at`
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW(); -- Simpan timestamp langsung dengan tipe TIMESTAMPTZ
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger untuk tabel `ingredients`
CREATE TRIGGER set_updated_at_ingredients
BEFORE UPDATE ON ingredients
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_ingredients_active;
DROP TRIGGER IF EXISTS set_updated_at_ingredients ON ingredients;
DROP TABLE IF EXISTS ingredients;
-- +goose StatementEnd
