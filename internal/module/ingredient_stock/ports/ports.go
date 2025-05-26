package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/ingredient_stock/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/ingredient_stock/entity"
	"github.com/jmoiron/sqlx"
)

type IngredientStockRepository interface {
	InsertNewIngredientStock(ctx context.Context, tx *sqlx.Tx, data *entity.IngredientStock) (*entity.IngredientStock, error)
	UpdateIngredientStock(ctx context.Context, tx *sqlx.Tx, data *entity.IngredientStock) (*entity.IngredientStock, error)
	FindListIngredientStock(ctx context.Context, limit, offset int, search string) ([]dto.GetListIngredientStock, int, error)
	FindIngredientStockByID(ctx context.Context, id int) (*entity.IngredientStock, error)
	SoftDeleteIngredientStockByID(ctx context.Context, tx *sqlx.Tx, id int) error
	FindIngredientStockByIDs(ctx context.Context, ids []int) ([]entity.IngredientStock, error)
	CountIngredientStockByID(ctx context.Context, id int) (int, error)
	DecrementStock(ctx context.Context, stockID string, amount int) error
}

type IngredientStockService interface {
	CreateNewIngredientStock(ctx context.Context, req *dto.CreateIngredientStockRequest) (*dto.CreateIngredientStockResponse, error)
	UpdateIngredientStock(ctx context.Context, req *dto.CreateIngredientStockRequest, id int) (*dto.CreateIngredientStockResponse, error)
	GetListIngredientStock(ctx context.Context, page, limit int, search string) (*dto.GetListIngredientStockResponse, error)
	GetDetailIngredientStock(ctx context.Context, id int) (*dto.GetListIngredientStock, error)
	RemoveIngredientStock(ctx context.Context, id int) error
}
