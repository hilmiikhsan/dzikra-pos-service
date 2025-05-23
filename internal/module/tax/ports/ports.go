package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/tax/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/tax/entity"
	"github.com/jmoiron/sqlx"
)

type TaxRepository interface {
	InsertNewTax(ctx context.Context, tx *sqlx.Tx, data *entity.Tax) (*entity.Tax, error)
	FindTax(ctx context.Context) (*entity.Tax, error)
	UpdateTax(ctx context.Context, tx *sqlx.Tx, data *entity.Tax) (*entity.Tax, error)
}

type TaxService interface {
	CreateOrUpdateTax(ctx context.Context, req *dto.CreateOrUpdateTaxRequest) (*dto.CreateOrUpdateTaxResponse, error)
}
