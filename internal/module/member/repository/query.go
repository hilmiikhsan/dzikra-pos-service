package repository

const (
	queryInsertNewMember = `
		INSERT INTO members
		(
			id,
			name,
			phone_number,
			email
		) VALUES (?, ?, ?, ?)
		RETURNING id, name, phone_number, email, created_at
	`

	queryFindListMember = `
        SELECT
            id,
            name,
            email,
            phone_number,
            created_at
        FROM members
        WHERE deleted_at IS NULL
          AND (
            name    ILIKE '%' || ? || '%'
            OR email     ILIKE '%' || ? || '%'
            OR phone_number ILIKE '%' || ? || '%'
          )
        ORDER BY created_at DESC, id DESC
        LIMIT ? OFFSET ?`

	queryCountFindListMember = `
        SELECT COUNT(*)
        FROM members
        WHERE deleted_at IS NULL
          AND (
            name    ILIKE '%' || ? || '%'
            OR email     ILIKE '%' || ? || '%'
            OR phone_number ILIKE '%' || ? || '%'
          )`

	queryUpdateMember = `
		UPDATE members
		SET 
			name = ?, 
			phone_number = ?, 
			email = ?
		WHERE id = ? AND deleted_at IS NULL
		RETURNING id, name, phone_number, email, created_at
	`

	queryFindMemberByID = `
    SELECT
      id,
      name,
      email,
      phone_number,
      created_at
    FROM members
    WHERE id = ? AND deleted_at IS NULL
  `

	querySoftDeleteMemberByID = `
    UPDATE members
    SET deleted_at = NOW()
    WHERE id = ? AND deleted_at IS NULL
  `

	queryFindMemberByEmailOrPhoneNumber = `
    SELECT
      id,
      name,
      email,
      phone_number,
      created_at
    FROM members
    WHERE 
      deleted_at IS NULL
      AND (email = $1 OR phone_number = $1)
  `
)
