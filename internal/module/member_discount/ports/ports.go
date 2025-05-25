package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member_discount/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member_discount/entity"
	"github.com/jmoiron/sqlx"
)

type MemberDiscountRepository interface {
	InsertNewMemberDiscount(ctx context.Context, tx *sqlx.Tx, data *entity.MemberDiscount) (*entity.MemberDiscount, error)
	FindDiscount(ctx context.Context) (*entity.MemberDiscount, error)
	UpdateMemberDiscount(ctx context.Context, tx *sqlx.Tx, data *entity.MemberDiscount) (*entity.MemberDiscount, error)
	FindFirstMemberDiscount(ctx context.Context) (*entity.MemberDiscount, error)
}

type MemberDiscountService interface {
	CreateOrUpdateMemberDiscount(ctx context.Context, req *dto.CreateOrUpdateMemberDiscountRequest) (*dto.CreateOrUpdateMemberDiscountResponse, error)
	CheckMemberDiscount(ctx context.Context, req *dto.CheckMemberDiscountRequest) (*dto.CheckMemberDiscountResponse, error)
}
