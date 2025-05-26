package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product/entity"
	"github.com/jmoiron/sqlx"
)

type ProductRepository interface {
	InsertNewProduct(ctx context.Context, tx *sqlx.Tx, data *entity.Product) (*entity.Product, error)
	UpdateProduct(ctx context.Context, tx *sqlx.Tx, data *entity.Product) (*entity.Product, error)
	CountProductByID(ctx context.Context, id int) (int, error)
	FindProductByID(ctx context.Context, id int) (*entity.Product, error)
	FindListProduct(ctx context.Context, limit, offset int, search string, productCategoryID int) ([]dto.GetListProduct, int, error)
	SoftDeleteProductByID(ctx context.Context, tx *sqlx.Tx, id int) error
	FindListProductRecipe(ctx context.Context, limit, offset int, search string) ([]dto.GetListProductRecipe, int, error)
	UpdateProductStock(ctx context.Context, tx *sqlx.Tx, data *entity.Product) error
	FindProductRecipeByProductIDs(ctx context.Context, productIDs []string) ([]entity.Product, error)
	UpdateProductsStock(ctx context.Context, stockMap map[string]int) error
}

type ProductService interface {
	CreateProduct(ctx context.Context, payloadFile dto.UploadFileRequest, req *dto.CreateOrUpdateProductRequest) (*dto.CreateOrUpdateProductResponse, error)
	UpdateProduct(ctx context.Context, payloadFile *dto.UploadFileRequest, req *dto.CreateOrUpdateProductRequest, id int) (*dto.CreateOrUpdateProductResponse, error)
	GetListProduct(ctx context.Context, page, limit int, search string, productCategoryID int) (*dto.GetListProductResponse, error)
	GetDetailProduct(ctx context.Context, id int) (*dto.GetListProduct, error)
	RemoveProduct(ctx context.Context, id int) error
}
