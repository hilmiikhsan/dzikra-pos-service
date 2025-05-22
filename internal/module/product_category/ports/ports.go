package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product_category/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product_category/entity"
	"github.com/jmoiron/sqlx"
)

type ProductCategoryRepository interface {
	InsertNewProductCategory(ctx context.Context, tx *sqlx.Tx, data *entity.ProductCategory) (*entity.ProductCategory, error)
	UpdateProductCategory(ctx context.Context, tx *sqlx.Tx, data *entity.ProductCategory) (*entity.ProductCategory, error)
	FindProductCategoryByID(ctx context.Context, id int) (*entity.ProductCategory, error)
	FindListProductCategory(ctx context.Context, limit, offset int, search string) ([]dto.GetProductCategoryResponse, int, error)
	SoftDeleteProductCategoryByID(ctx context.Context, tx *sqlx.Tx, id int) error
}

type ProductCategoryService interface {
	CreateOrUpdateProductCategory(ctx context.Context, req *dto.CreateOrUpdateProductCategoryRequest) (*dto.CreateOrUpdateProductCategoryResponse, error)
	UpdateProductCategory(ctx context.Context, req *dto.CreateOrUpdateProductCategoryRequest, id int) (*dto.CreateOrUpdateProductCategoryResponse, error)
	GetDetailProductCategory(ctx context.Context, id int) (*dto.GetProductCategoryResponse, error)
	GetListProductCategory(ctx context.Context, page, limit int, search string) (*dto.GetListProductCategoryResponse, error)
	RemoveProductCategory(ctx context.Context, id int) error
}
