package rest

import (
	externalUser "github.com/Digitalkeun-Creative/be-dzikra-pos-service/external/user"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/adapter"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/middleware"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member/ports"
	memberRepository "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member/repository"
	memberService "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member/service"
)

type memberHandler struct {
	service    ports.MemberService
	middleware middleware.AuthMiddleware
	validator  adapter.Validator
}

func NewmemberHandler() *memberHandler {
	var handler = new(memberHandler)

	// validator
	validator := adapter.Adapters.Validator

	// external
	externalAuth := &externalUser.External{}

	// middleware
	middlewareHandler := middleware.NewAuthMiddleware(externalAuth)

	// repository
	memberRepository := memberRepository.NewMemberRepository(adapter.Adapters.DzikraPostgres)

	// faq service
	faqService := memberService.NewMemberService(
		adapter.Adapters.DzikraPostgres,
		memberRepository,
	)

	// handler
	handler.service = faqService
	handler.middleware = *middlewareHandler
	handler.validator = validator

	return handler
}
