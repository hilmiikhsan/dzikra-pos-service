package rest

import (
	externalUser "github.com/Digitalkeun-Creative/be-dzikra-pos-service/external/user"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/adapter"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/middleware"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/dashboard/ports"
	dashboardService "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/dashboard/service"
	expensesRepository "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/expenses/repository"
	transactionHistoryRepository "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/transaction_history/repository"
)

type dashboardHandler struct {
	service    ports.DashboardService
	middleware middleware.AuthMiddleware
	validator  adapter.Validator
}

func NewDashboardHandler() *dashboardHandler {
	var handler = new(dashboardHandler)

	// validator
	validator := adapter.Adapters.Validator

	// external
	externalAuth := &externalUser.External{}

	// middleware
	middlewareHandler := middleware.NewAuthMiddleware(externalAuth)

	// repository
	transactionHistoryRepository := transactionHistoryRepository.NewTransactionHistoryRepository(adapter.Adapters.DzikraPostgres)
	expensesRepository := expensesRepository.NewExpensesRepository(adapter.Adapters.DzikraPostgres)

	// dashboard service
	dashboardService := dashboardService.NewDashboardService(
		adapter.Adapters.DzikraPostgres,
		transactionHistoryRepository,
		expensesRepository,
	)

	// handler
	handler.service = dashboardService
	handler.middleware = *middlewareHandler
	handler.validator = validator

	return handler
}
