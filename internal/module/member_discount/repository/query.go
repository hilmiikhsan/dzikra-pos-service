package repository

const (
	queryInsertNewMemberDiscount = `
		INSERT INTO member_discounts 
		(
			discount, 
			updated_at
		) VALUES (?, NOW())
		RETURNING id, discount, updated_at
	`

	queryFindMemberDiscount = `
		SELECT
			id,
			discount,
			updated_at
		FROM member_discounts
		WHERE deleted_at IS NULL
		LIMIT 1
	`

	queryUpdateMemberDiscount = `
		UPDATE member_discounts
		SET
			discount = ?,
			updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
		RETURNING id, discount, updated_at
	`

	queryGetFirstMemberDiscount = `
		SELECT discount
		FROM member_discounts
		ORDER BY created_at DESC
		LIMIT 1
	`
)
