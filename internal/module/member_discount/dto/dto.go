package dto

import (
	member "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member/dto"
)

type CreateOrUpdateMemberDiscountRequest struct {
	Discount int `json:"discount" validate:"required,numeric,number"`
}

type CreateOrUpdateMemberDiscountResponse struct {
	ID        int    `json:"id"`
	Discount  string `json:"discount"`
	UpdatedAt string `json:"updated_at"`
}

type CheckMemberDiscountRequest struct {
	Identifier string `json:"identifier" validate:"required,min=2,max=20,xss_safe"`
}

type CheckMemberDiscountResponse struct {
	Discount string               `json:"discount"`
	Member   member.GetListMember `json:"member"`
}
