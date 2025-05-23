package dto

type CreateOrUpdateTaxRequest struct {
	TaxAmount int `json:"tax_amount" validate:"required,numeric,number"`
}

type CreateOrUpdateTaxResponse struct {
	ID  int    `json:"id"`
	Tax string `json:"tax"`
}
