package service

import (
	memberPorts "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member/ports"
	"github.com/jmoiron/sqlx"
)

var _ memberPorts.MemberService = &memberService{}

type memberService struct {
	db               *sqlx.DB
	memberRepository memberPorts.MemberRepository
}

func NewMemberService(
	db *sqlx.DB,
	memberRepository memberPorts.MemberRepository,
) *memberService {
	return &memberService{
		db:               db,
		memberRepository: memberRepository,
	}
}
