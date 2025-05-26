package entity

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID                       uuid.UUID `db:"id"`
	Status                   string    `db:"status"`
	PhoneNumber              string    `db:"phone_number"`
	Name                     string    `db:"name"`
	Email                    string    `db:"email"`
	IsMember                 bool      `db:"is_member"`
	TotalQuantity            int       `db:"total_quantity"`
	TotalProductAmount       int       `db:"total_product_amount"`
	TotalAmount              int       `db:"total_amount"`
	VPaymentID               string    `db:"v_payment_id"`
	VPaymentRedirectUrl      string    `db:"v_payment_redirect_url"`
	VTransactionID           string    `db:"v_transaction_id"`
	DiscountPercentage       int       `db:"discount_percentage"`
	ChangeMoney              int       `db:"change_money"`
	PaymentType              string    `db:"payment_type"`
	TotalMoney               int       `db:"total_money"`
	TableNumber              int       `db:"table_number"`
	TotalProductCapitalPrice int       `db:"total_product_capital_price"`
	TaxAmount                int       `db:"tax_amount"`
	Notes                    string    `db:"notes"`
	CreatedAt                time.Time `db:"created_at"`
}
