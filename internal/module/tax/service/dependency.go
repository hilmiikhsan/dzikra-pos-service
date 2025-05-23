package service

import (
	taxPorts "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/tax/ports"
	"github.com/jmoiron/sqlx"
)

var _ taxPorts.TaxService = &taxService{}

type taxService struct {
	db            *sqlx.DB
	taxRepository taxPorts.TaxRepository
}

func NewTaxService(
	db *sqlx.DB,
	taxRepository taxPorts.TaxRepository,
) *taxService {
	return &taxService{
		db:            db,
		taxRepository: taxRepository,
	}
}
