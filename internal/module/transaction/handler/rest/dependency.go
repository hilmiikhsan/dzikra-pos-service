package rest

import (
	externalTransaction "github.com/Digitalkeun-Creative/be-dzikra-pos-service/external/transaction"
	externalUser "github.com/Digitalkeun-Creative/be-dzikra-pos-service/external/user"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/adapter"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/middleware"
	ingredientRepository "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/ingredient/repository"
	ingredientStockRepository "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/ingredient_stock/repository"
	memberRepository "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member/repository"
	memberDiscountRepository "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member_discount/repository"
	productRepository "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product/repository"
	recipeRepository "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/recipe/repository"
	recipeService "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/recipe/service"
	taxRepository "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/tax/repository"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/transaction/ports"
	transactionRepository "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/transaction/repository"
	transactionService "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/transaction/service"
	transactionItemRepository "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/transaction_item/repository"
)

type transactionHandler struct {
	service    ports.TransactionService
	middleware middleware.AuthMiddleware
	validator  adapter.Validator
}

func NewTransactionHandler() *transactionHandler {
	var handler = new(transactionHandler)

	// validator
	validator := adapter.Adapters.Validator

	// external
	externalAuth := &externalUser.External{}
	externalTransaction := &externalTransaction.External{}

	// middleware
	middlewareHandler := middleware.NewAuthMiddleware(externalAuth)

	// repository
	transactionRepository := transactionRepository.NewTransactionRepository(adapter.Adapters.DzikraPostgres)
	transactionItemRepository := transactionItemRepository.NewTransactionItemRepository(adapter.Adapters.DzikraPostgres)
	memberRepository := memberRepository.NewMemberRepository(adapter.Adapters.DzikraPostgres)
	productRepository := productRepository.NewProductRepository(adapter.Adapters.DzikraPostgres)
	memberDiscountRepository := memberDiscountRepository.NewMemberDiscountRepository(adapter.Adapters.DzikraPostgres)
	taxRepository := taxRepository.NewTaxRepository(adapter.Adapters.DzikraPostgres)
	recipeRepository := recipeRepository.NewRecipeRepository(adapter.Adapters.DzikraPostgres)
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

	// transaction service
	transactionService := transactionService.NewTransactionService(
		adapter.Adapters.DzikraPostgres,
		transactionRepository,
		transactionItemRepository,
		memberRepository,
		productRepository,
		memberDiscountRepository,
		taxRepository,
		recipeService,
		externalTransaction,
	)

	// handler
	handler.service = transactionService
	handler.middleware = *middlewareHandler
	handler.validator = validator

	return handler
}
