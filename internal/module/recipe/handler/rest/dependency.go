package rest

import (
	externalUser "github.com/Digitalkeun-Creative/be-dzikra-pos-service/external/user"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/adapter"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/middleware"
	ingredientRepository "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/ingredient/repository"
	ingredientStockRepository "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/ingredient_stock/repository"
	productRepository "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product/repository"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/recipe/ports"
	recipeRepository "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/recipe/repository"
	recipeService "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/recipe/service"
)

type recipeHandler struct {
	service    ports.RecipeService
	middleware middleware.AuthMiddleware
	validator  adapter.Validator
}

func NewRecipeHandler() *recipeHandler {
	var handler = new(recipeHandler)

	// validator
	validator := adapter.Adapters.Validator

	// external
	externalAuth := &externalUser.External{}

	// middleware
	middlewareHandler := middleware.NewAuthMiddleware(externalAuth)

	// repository
	recipeRepository := recipeRepository.NewRecipeRepository(adapter.Adapters.DzikraPostgres)
	productRepository := productRepository.NewProductRepository(adapter.Adapters.DzikraPostgres)
	ingredientStockRepository := ingredientStockRepository.NewIngredientStockRepository(adapter.Adapters.DzikraPostgres)
	ingredientRepository := ingredientRepository.NewIngredientRepository(adapter.Adapters.DzikraPostgres)

	// recipe service
	recipeService := recipeService.NewRecipeService(
		adapter.Adapters.DzikraPostgres,
		recipeRepository,
		productRepository,
		ingredientStockRepository,
		ingredientRepository,
	)

	// handler
	handler.service = recipeService
	handler.middleware = *middlewareHandler
	handler.validator = validator

	return handler
}
