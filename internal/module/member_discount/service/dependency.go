package service

import (
	memberDiscountPorts "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member_discount/ports"
	"github.com/jmoiron/sqlx"
)

var _ memberDiscountPorts.MemberDiscountService = &memberDiscountService{}

type memberDiscountService struct {
	db                       *sqlx.DB
	memberDiscountRepository memberDiscountPorts.MemberDiscountRepository
}

func NewMemberDiscountService(
	db *sqlx.DB,
	memberDiscountRepository memberDiscountPorts.MemberDiscountRepository,
) *memberDiscountService {
	return &memberDiscountService{
		db:                       db,
		memberDiscountRepository: memberDiscountRepository,
	}
}
