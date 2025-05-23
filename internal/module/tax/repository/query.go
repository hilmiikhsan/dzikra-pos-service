package repository

const (
	queryInsertNewTax = `
		INSERT INTO tax (tax) VALUES(?) RETURNING id, tax
	`

	queryFindTax = `
		SELECT
			id,
			tax
		FROM tax
		WHERE deleted_at IS NULL
		LIMIT 1
	`

	queryUpdateTax = `
		UPDATE tax
		SET
			tax = ?,
			updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
		RETURNING id, tax
	`
)
