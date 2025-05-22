package rest

import (
	externalUser "github.com/Digitalkeun-Creative/be-dzikra-pos-service/external/user"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/adapter"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/middleware"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product_category/ports"
	productCategoryrepository "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product_category/repository"
	productCategoryService "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product_category/service"
)

type productCategoryHandler struct {
	service    ports.ProductCategoryService
	middleware middleware.AuthMiddleware
	validator  adapter.Validator
}

func NewProductCategoryHandler() *productCategoryHandler {
	var handler = new(productCategoryHandler)

	// validator
	validator := adapter.Adapters.Validator

	// external
	externalAuth := &externalUser.External{}

	// middleware
	middlewareHandler := middleware.NewAuthMiddleware(externalAuth)

	// repository
	productCategoryrepository := productCategoryrepository.NewProductCategoryRepository(adapter.Adapters.DzikraPostgres)

	// member discount service
	productCategoryService := productCategoryService.NewProductCategoryService(
		adapter.Adapters.DzikraPostgres,
		productCategoryrepository,
	)

	// handler
	handler.service = productCategoryService
	handler.middleware = *middlewareHandler
	handler.validator = validator

	return handler
}
