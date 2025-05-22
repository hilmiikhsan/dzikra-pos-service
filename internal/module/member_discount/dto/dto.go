package dto

type CreateOrUpdateMemberDiscountRequest struct {
	Discount int `json:"discount" validate:"required,numeric,number"`
}

type CreateOrUpdateMemberDiscountResponse struct {
	ID        int    `json:"id"`
	Discount  string `json:"discount"`
	UpdatedAt string `json:"updated_at"`
}
