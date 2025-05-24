package rest

import (
	externalUser "github.com/Digitalkeun-Creative/be-dzikra-pos-service/external/user"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/adapter"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/middleware"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/ingredient_stock/ports"
	ingredientStockRepository "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/ingredient_stock/repository"
	ingredientStockService "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/ingredient_stock/service"
)

type ingredientStockHandler struct {
	service    ports.IngredientStockService
	middleware middleware.AuthMiddleware
	validator  adapter.Validator
}

func NewIngredientStockHandler() *ingredientStockHandler {
	var handler = new(ingredientStockHandler)

	// validator
	validator := adapter.Adapters.Validator

	// external
	externalAuth := &externalUser.External{}

	// middleware
	middlewareHandler := middleware.NewAuthMiddleware(externalAuth)

	// repository
	ingredientStockRepository := ingredientStockRepository.NewIngredientStockRepository(adapter.Adapters.DzikraPostgres)

	// ingredient stock service
	ingredientStockService := ingredientStockService.NewIngredientStockService(
		adapter.Adapters.DzikraPostgres,
		ingredientStockRepository,
	)

	// handler
	handler.service = ingredientStockService
	handler.middleware = *middlewareHandler
	handler.validator = validator

	return handler
}
