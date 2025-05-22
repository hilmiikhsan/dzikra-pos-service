package entity

import (
	"time"

	"github.com/google/uuid"
)

type Member struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	PhoneNumber string    `db:"phone_number"`
	Email       string    `db:"email"`
	CreatedAt   time.Time `db:"created_at"`
}
