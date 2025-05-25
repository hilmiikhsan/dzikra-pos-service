package rest

import (
	externalUser "github.com/Digitalkeun-Creative/be-dzikra-pos-service/external/user"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/adapter"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/middleware"
	memberRepository "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member/repository"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member_discount/ports"
	memberDiscountrepository "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member_discount/repository"
	memberDiscountService "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member_discount/service"
)

type memberDiscountHandler struct {
	service    ports.MemberDiscountService
	middleware middleware.AuthMiddleware
	validator  adapter.Validator
}

func NewMemberDiscountHandler() *memberDiscountHandler {
	var handler = new(memberDiscountHandler)

	// validator
	validator := adapter.Adapters.Validator

	// external
	externalAuth := &externalUser.External{}

	// middleware
	middlewareHandler := middleware.NewAuthMiddleware(externalAuth)

	// repository
	memberDiscountrepository := memberDiscountrepository.NewMemberDiscountRepository(adapter.Adapters.DzikraPostgres)
	memberRepository := memberRepository.NewMemberRepository(adapter.Adapters.DzikraPostgres)

	// member discount service
	memberDiscount := memberDiscountService.NewMemberDiscountService(
		adapter.Adapters.DzikraPostgres,
		memberDiscountrepository,
		memberRepository,
	)

	// handler
	handler.service = memberDiscount
	handler.middleware = *middlewareHandler
	handler.validator = validator

	return handler
}
