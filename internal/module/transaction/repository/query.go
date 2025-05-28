package repository

const (
	queryInsertNewTransaction = `
		INSERT INTO transactions
		(
			id,
			status,
			phone_number,
			name,
			email,
			is_member,
			total_quantity,
			total_product_amount,
			total_amount,
			v_payment_id,
			v_payment_redirect_url,
			v_transaction_id,
			discount_percentage,
			change_money,
			payment_type,
			total_money,
			table_number,
			total_product_capital_price,
			tax_amount,
			notes,
			created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	queryUpdateCashField = `
		UPDATE transactions
		SET total_money = $1, change_money = $2
		WHERE id = $3
	`

	queryUpdateTransactionByID = `
		UPDATE transactions
        SET 
			v_transaction_id = $1,
            v_payment_id = $2,
            v_payment_redirect_url = $3
        WHERE id = $4
	`

	queryFindListTransaction = `
		SELECT
			id,
			status,
			phone_number,
			name,
			email,
			is_member,
			total_quantity,
			total_product_amount,
			total_product_capital_price,
			total_amount,
			discount_percentage,
			v_transaction_id,
			v_payment_id,
			v_payment_redirect_url,
			payment_type,
			table_number,
			notes,
			tax_amount,
			created_at
		FROM transactions
		WHERE status IN (?, ?)
			AND (
				name ILIKE '%' || ? || '%'
				OR email ILIKE '%' || ? || '%'
				OR phone_number ILIKE '%' || ? || '%'
			)
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	queryCountFindListTransaction = `
		SELECT COUNT(*) 
		FROM transactions
		WHERE status IN (?, ?)
		AND (
			name ILIKE '%' || ? || '%'
			OR email ILIKE '%' || ? || '%'
			OR phone_number ILIKE '%' || ? || '%'
		)
	`

	queryFindTransactionWithItemsByID = `
		SELECT
			id, 
			status, 
			phone_number, 
			name, 
			email, 
			is_member,
			total_quantity, 
			total_product_amount, 
			total_amount,
			v_payment_id, 
			v_payment_redirect_url,
			v_transaction_id,
			discount_percentage, 
			change_money, 
			payment_type, 
			total_money,
			table_number, 
			total_product_capital_price, 
			tax_amount,
			notes, created_at
		FROM transactions
		WHERE id = ?
	`

	queryFindItemsByTransactionID = `
		SELECT
			id, 
			transaction_id, 
			product_id, 
			quantity, 
			total_amount,
			product_name, 
			product_price, 
			product_capital_price,
			total_amount_capital_price
		FROM transaction_items
		WHERE transaction_id = ?
		ORDER BY id
	`

	queryFindTransactionByVPaymentID = `
		SELECT
			id, 
			transaction_id, 
			product_id, 
			quantity, 
			total_amount,
			product_name, 
			product_price, 
			product_capital_price,
			total_amount_capital_price
		FROM transactions
		WHERE v_payment_id = ? AND deleted_at IS NULL
	`

	queryUpdateTransactionStatus = `
		UPDATE transactions 
		SET 
			status = ?, 
			updated_at = NOW()
		WHERE id = ?
		RETURNING id, status, phone_number, name, email
	`

	queryDuplicateToTransactionHistory = `
		INSERT INTO transaction_history
		SELECT * FROM transactions WHERE id = ?
	`

	queryDuplicateToTransactionItemsHistory = `
		INSERT INTO transaction_item_history
		SELECT * FROM transaction_items WHERE transaction_id = ?
	`
)
