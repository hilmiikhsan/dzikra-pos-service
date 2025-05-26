package repository

const (
	queryInsertNewExpenses = `
		INSERT INTO expenses 
		(
			name, 
			cost, 
			date
		) VALUES (?, ?, ?)
		RETURNING id, name, cost, date, created_at
	`

	queryFindListExpenses = `
		SELECT
			id,
			name,
			cost,
			date,
			created_at,
			updated_at
		FROM expenses
		WHERE deleted_at IS NULL
		AND (
			name ILIKE '%' || ? || '%'
			OR cost::TEXT ILIKE '%' || ? || '%'
		)
		ORDER BY created_at DESC, id DESC
		LIMIT ? OFFSET ?
	`

	queryCountFindListExpenses = `
		SELECT COUNT(*) 
		FROM expenses
		WHERE deleted_at IS NULL
		AND (
			name ILIKE '%' || ? || '%'
			OR cost::TEXT ILIKE '%' || ? || '%'
		)
	`

	queryFindExpensesByID = `
		SELECT
			id,
			name,
			cost,
			date,
			created_at,
			updated_at
		FROM expenses
		WHERE id = ? AND deleted_at IS NULL
	`

	queryUpdateExpenses = `
		UPDATE expenses
		SET
			name = ?,
			cost = ?,
			date = ?,
			updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
		RETURNING id, name, cost, date, created_at, updated_at
	`

	querySoftDeleteExpensesByID = `
		UPDATE expenses
		SET deleted_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`
)
