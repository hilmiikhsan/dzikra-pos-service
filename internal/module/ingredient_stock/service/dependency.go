package service

import (
	ingredientStockorts "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/ingredient_stock/ports"
	"github.com/jmoiron/sqlx"
)

var _ ingredientStockorts.IngredientStockService = &ingredientStockService{}

type ingredientStockService struct {
	db                        *sqlx.DB
	ingredientStockRepository ingredientStockorts.IngredientStockRepository
}

func NewIngredientStockService(
	db *sqlx.DB,
	ingredientStockRepository ingredientStockorts.IngredientStockRepository,
) *ingredientStockService {
	return &ingredientStockService{
		db:                        db,
		ingredientStockRepository: ingredientStockRepository,
	}
}
