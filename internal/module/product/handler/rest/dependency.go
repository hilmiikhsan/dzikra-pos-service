package rest

import (
	externalUser "github.com/Digitalkeun-Creative/be-dzikra-pos-service/external/user"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/adapter"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/infrastructure/config"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/integration/storage/minio"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/middleware"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product/ports"
	productRepository "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product/repository"
	productService "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product/service"
	productCategoryRepository "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product_category/repository"
	recipeRepository "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/recipe/repository"
)

type productHandler struct {
	service    ports.ProductService
	middleware middleware.AuthMiddleware
	validator  adapter.Validator
}

func NewProductHandler() *productHandler {
	var handler = new(productHandler)

	// validator
	validator := adapter.Adapters.Validator

	// external
	externalAuth := &externalUser.External{}

	// minio service
	minioService := minio.NewMinioService(adapter.Adapters.DzikraMinio, config.Envs.MinioStorage.Bucket)

	// middleware
	middlewareHandler := middleware.NewAuthMiddleware(externalAuth)

	// repository
	productRepository := productRepository.NewProductRepository(adapter.Adapters.DzikraPostgres)
	productCategoryRepository := productCategoryRepository.NewProductCategoryRepository(adapter.Adapters.DzikraPostgres)
	recipeRepository := recipeRepository.NewRecipeRepository(adapter.Adapters.DzikraPostgres)

	// product service
	productService := productService.NewProductService(
		adapter.Adapters.DzikraPostgres,
		productRepository,
		productCategoryRepository,
		recipeRepository,
		minioService,
	)

	// handler
	handler.service = productService
	handler.middleware = *middlewareHandler
	handler.validator = validator

	return handler
}
