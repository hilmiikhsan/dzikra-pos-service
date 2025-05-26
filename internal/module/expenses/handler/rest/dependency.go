package rest

import (
	externalUser "github.com/Digitalkeun-Creative/be-dzikra-pos-service/external/user"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/adapter"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/middleware"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/expenses/ports"
	expensesRepository "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/expenses/repository"
	expensesService "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/expenses/service"
)

type expensesHandler struct {
	service    ports.ExpensesService
	middleware middleware.AuthMiddleware
	validator  adapter.Validator
}

func NewExpensesHandler() *expensesHandler {
	var handler = new(expensesHandler)

	// validator
	validator := adapter.Adapters.Validator

	// external
	externalAuth := &externalUser.External{}

	// middleware
	middlewareHandler := middleware.NewAuthMiddleware(externalAuth)

	// repository
	expensesRepository := expensesRepository.NewExpensesRepository(adapter.Adapters.DzikraPostgres)

	// expenses service
	expensesService := expensesService.NewExpensesService(
		adapter.Adapters.DzikraPostgres,
		expensesRepository,
	)

	// handler
	handler.service = expensesService
	handler.middleware = *middlewareHandler
	handler.validator = validator

	return handler
}
