package service

import (
	ingredientPorts "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/ingredient/ports"
	ingredientStockPorts "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/ingredient_stock/ports"
	productPorts "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product/ports"
	recipePorts "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/recipe/ports"
	"github.com/jmoiron/sqlx"
)

var _ recipePorts.RecipeService = &recipeService{}

type recipeService struct {
	db                        *sqlx.DB
	recipeRepository          recipePorts.RecipeRepository
	productRepository         productPorts.ProductRepository
	ingredientStockRepository ingredientStockPorts.IngredientStockRepository
	ingredientRepository      ingredientPorts.IngredientRepository
}

func NewRecipeService(
	db *sqlx.DB,
	recipeRepository recipePorts.RecipeRepository,
	productRepository productPorts.ProductRepository,
	ingredientStockRepository ingredientStockPorts.IngredientStockRepository,
	ingredientRepository ingredientPorts.IngredientRepository,
) *recipeService {
	return &recipeService{
		db:                        db,
		recipeRepository:          recipeRepository,
		productRepository:         productRepository,
		ingredientStockRepository: ingredientStockRepository,
		ingredientRepository:      ingredientRepository,
	}
}
