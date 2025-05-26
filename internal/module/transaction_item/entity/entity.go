package entity

import (
	"time"

	"github.com/google/uuid"
)

type TransactionItem struct {
	ID                      int       `db:"id"`
	Quantity                int       `db:"quantity"`
	TotalAmount             int       `db:"total_amount"`
	ProductName             string    `db:"product_name"`
	ProductPrice            int       `db:"product_price"`
	TransactionID           uuid.UUID `db:"transaction_id"`
	ProductID               int       `db:"product_id"`
	ProductCapitalPrice     int       `db:"product_capital_price"`
	TotalAmountCapitalPrice int       `db:"total_amount_capital_price"`
	CreatedAt               time.Time `db:"created_at"`
}
