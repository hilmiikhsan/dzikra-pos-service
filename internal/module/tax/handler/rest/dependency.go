package rest

import (
	externalUser "github.com/Digitalkeun-Creative/be-dzikra-pos-service/external/user"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/adapter"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/middleware"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/tax/ports"
	taxRepository "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/tax/repository"
	taxService "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/tax/service"
)

type taxHandler struct {
	service    ports.TaxService
	middleware middleware.AuthMiddleware
	validator  adapter.Validator
}

func NewTaxHandler() *taxHandler {
	var handler = new(taxHandler)

	// validator
	validator := adapter.Adapters.Validator

	// external
	externalAuth := &externalUser.External{}

	// middleware
	middlewareHandler := middleware.NewAuthMiddleware(externalAuth)

	// repository
	taxRepository := taxRepository.NewTaxRepository(adapter.Adapters.DzikraPostgres)

	// tax service
	taxService := taxService.NewTaxService(
		adapter.Adapters.DzikraPostgres,
		taxRepository,
	)

	// handler
	handler.service = taxService
	handler.middleware = *middlewareHandler
	handler.validator = validator

	return handler
}
