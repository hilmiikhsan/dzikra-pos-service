package entity

import "time"

type MemberDiscount struct {
	ID        int       `db:"id"`
	Discount  int       `db:"discount"`
	UpdatedAt time.Time `db:"updated_at"`
}
