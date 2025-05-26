package repository

const (
	queryInsertNewIngredientStock = `
		INSERT INTO ingredient_stocks
		(
			name,
			unit,
			price_per_amount_stock,
			required_stock,
			amount_price_required_stock,
			amount_stock_per_price
		) VALUES (?, ?, ?, ?, ?, ?) 
		RETURNING 
			id,
			name,
			unit,
			CAST(price_per_amount_stock           AS INTEGER) AS price_per_amount_stock,
			required_stock,
			CAST(amount_price_required_stock      AS INTEGER) AS amount_price_required_stock,
			CAST(amount_stock_per_price           AS INTEGER) AS amount_stock_per_price
	`

	queryUpdateIngredientStock = `
		UPDATE ingredient_stocks
		SET
			name = ?,
			unit = ?,
			price_per_amount_stock = ?,
			required_stock = ?,
			amount_price_required_stock = ?,
			amount_stock_per_price = ?
		WHERE id = ? AND deleted_at IS NULL
		RETURNING
			id,
			name,
			unit,
			CAST(price_per_amount_stock           AS INTEGER) AS price_per_amount_stock,
			required_stock,
			CAST(amount_price_required_stock      AS INTEGER) AS amount_price_required_stock,
			CAST(amount_stock_per_price           AS INTEGER) AS amount_stock_per_price
	`

	queryFindListIngredientStock = `
		SELECT
			id,
			name,
			unit,
			CAST(price_per_amount_stock           AS INTEGER) AS price_per_amount_stock,
			required_stock,
			CAST(amount_price_required_stock      AS INTEGER) AS amount_price_required_stock,
			CAST(amount_stock_per_price           AS INTEGER) AS amount_stock_per_price
		FROM ingredient_stocks
		WHERE deleted_at IS NULL
		AND (name ILIKE '%' || ? || '%')
		ORDER BY created_at DESC, id DESC
		LIMIT ? OFFSET ?
	`

	queryCountFindListIngredientStock = `
		SELECT COUNT(*)
		FROM ingredient_stocks
		WHERE deleted_at IS NULL
		AND (name ILIKE '%' || ? || '%')
	`

	queryFindIngredientStockByID = `
		SELECT
			id,
			name,
			unit,
			CAST(price_per_amount_stock           AS INTEGER) AS price_per_amount_stock,
			required_stock,
			CAST(amount_price_required_stock      AS INTEGER) AS amount_price_required_stock,
			CAST(amount_stock_per_price           AS INTEGER) AS amount_stock_per_price
		FROM ingredient_stocks
		WHERE id = ? AND deleted_at IS NULL
	`

	querySoftDeleteIngredientStockByID = `
		UPDATE ingredient_stocks
		SET deleted_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`

	queryFindIngredientStockByIDs = `
		SELECT id, name, required_stock, unit,
             CAST(price_per_amount_stock           AS INTEGER) AS price_per_amount_stock,
             CAST(amount_stock_per_price           AS INTEGER) AS amount_stock_per_price,
             CAST(amount_price_required_stock      AS INTEGER) AS amount_price_required_stock
      FROM ingredient_stocks
      WHERE id = ANY($1::int[]) AND deleted_at IS NULL
	`

	queryCountIngredientStockByID = `
		SELECT COUNT(*) FROM ingredient_stocks WHERE id = $1 AND deleted_at IS NULL
	`

	queryDecrementStock = `
		UPDATE ingredient_stocks SET required_stock = required_stock - ? WHERE id = ? AND deleted_at IS NULL
	`
)
