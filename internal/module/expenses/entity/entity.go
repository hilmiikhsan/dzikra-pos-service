package entity

import "time"

type Expenses struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Cost      int       `db:"cost"`
	Date      time.Time `db:"date"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
