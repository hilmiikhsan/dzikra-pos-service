package repository

const (
	queryInsertNewProductCategory = `
		INSERT INTO product_categories (name) VALUES (?) RETURNING id, name
	`

	queryUpdateProductCategory = `
		UPDATE product_categories SET name = ? 
		WHERE id = ? AND deleted_at IS NULL
		RETURNING id, name
	`

	queryFindProductCategoryByID = `
		SELECT id, name FROM product_categories
		WHERE id = ? AND deleted_at IS NULL
	`

	queryFindListproductCategory = `
		SELECT id, name FROM product_categories
		WHERE deleted_at IS NULL
		AND (name ILIKE '%' || ? || '%')
		ORDER BY created_at DESC, id DESC
		LIMIT ? OFFSET ?
	`

	queryCountFindListproductCategory = `
		SELECT COUNT(*) FROM product_categories
		WHERE deleted_at IS NULL
		AND (name ILIKE '%' || ? || '%')
	`

	querySoftDeleteProductCategoryByID = `
		UPDATE product_categories SET deleted_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`
)
