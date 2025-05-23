package service

import (
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/integration/storage/minio"
	productPorts "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product/ports"
	productCategoryPorts "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product_category/ports"
	recipePorts "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/recipe/ports"
	"github.com/jmoiron/sqlx"
)

var _ productPorts.ProductService = &productService{}

type productService struct {
	db                        *sqlx.DB
	productRepository         productPorts.ProductRepository
	productCategoryRepository productCategoryPorts.ProductCategoryRepository
	recipeRepository          recipePorts.RecipeRepository
	minioService              minio.MinioService
}

func NewProductService(
	db *sqlx.DB,
	productRepository productPorts.ProductRepository,
	productCategoryRepository productCategoryPorts.ProductCategoryRepository,
	recipeRepository recipePorts.RecipeRepository,
	minioService minio.MinioService,
) *productService {
	return &productService{
		db:                        db,
		productRepository:         productRepository,
		productCategoryRepository: productCategoryRepository,
		recipeRepository:          recipeRepository,
		minioService:              minioService,
	}
}
