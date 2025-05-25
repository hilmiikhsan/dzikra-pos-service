package service

import (
	memberPorts "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member/ports"
	memberDiscountPorts "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member_discount/ports"
	"github.com/jmoiron/sqlx"
)

var _ memberDiscountPorts.MemberDiscountService = &memberDiscountService{}

type memberDiscountService struct {
	db                       *sqlx.DB
	memberDiscountRepository memberDiscountPorts.MemberDiscountRepository
	memberRepository         memberPorts.MemberRepository
}

func NewMemberDiscountService(
	db *sqlx.DB,
	memberDiscountRepository memberDiscountPorts.MemberDiscountRepository,
	memberRepository memberPorts.MemberRepository,
) *memberDiscountService {
	return &memberDiscountService{
		db:                       db,
		memberDiscountRepository: memberDiscountRepository,
		memberRepository:         memberRepository,
	}
}
