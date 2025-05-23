package entity

type Recipe struct {
	ID           int `db:"id"`
	CapitalPrice int `db:"capital_price"`
}
