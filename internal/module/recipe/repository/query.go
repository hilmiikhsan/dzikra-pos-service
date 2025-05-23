package repository

const (
	queryInsertNewRecipe = `
		INSERT INTO recipes (capital_price) VALUES (?) RETURNING id, capital_price
	`
)
