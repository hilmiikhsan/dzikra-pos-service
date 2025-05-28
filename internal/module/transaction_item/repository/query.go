package repository

const (
	queryInsertNewTransactionItem = `
		INSERT INTO transaction_items
		(
			quantity,
			total_amount,
			product_name,
			product_price,
			transaction_id,
			product_id,
			product_capital_price,
			total_amount_capital_price
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	queryFindTransactionItemByTransactionID = `
		SELECT
			id,
			quantity,
			total_amount,
			product_name,
			product_price,
			transaction_id,
			product_id,
			product_capital_price,
			total_amount_capital_price
		FROM transaction_items
		WHERE transaction_id = ?
		AND deleted_at IS NULL
	`
)
