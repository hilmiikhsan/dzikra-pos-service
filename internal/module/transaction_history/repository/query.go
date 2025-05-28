package repository

const (
	querySumTotalAmount = `
		SELECT COALESCE(SUM(total_amount),0) 
		FROM transaction_histories
		WHERE created_at BETWEEN $1 AND $2
	`

	queryCountTransactionHistory = `
		SELECT COUNT(*) 
		FROM transaction_histories
		WHERE created_at BETWEEN $1 AND $2
	`

	querySumTotalQuantity = `
		SELECT COALESCE(SUM(total_quantity),0) 
		FROM transaction_histories
		WHERE created_at BETWEEN $1 AND $2
	`

	querySumTotalCapital = `
		SELECT COALESCE(SUM(total_product_capital_price),0) 
		FROM transaction_histories
		WHERE created_at BETWEEN $1 AND $2
	`
)
