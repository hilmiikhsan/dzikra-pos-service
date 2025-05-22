package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member/entity"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type MemberRepository interface {
	InsertNewMember(ctx context.Context, tx *sqlx.Tx, data *entity.Member) (*entity.Member, error)
	FindListMember(ctx context.Context, limit, offset int, search string) ([]dto.GetListMember, int, error)
	UpdateMember(ctx context.Context, tx *sqlx.Tx, data *entity.Member, id uuid.UUID) (*entity.Member, error)
	FindMemberByID(ctx context.Context, id uuid.UUID) (*entity.Member, error)
	SoftDeleteMemberByID(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) error
}

type MemberService interface {
	CreateMember(ctx context.Context, req *dto.CreateOrUpdateMemberRequest) (*dto.CreateOrUpdateMemberResponse, error)
	GetListMember(ctx context.Context, page, limit int, search string) (*dto.GetListMemberResponse, error)
	UpdateMember(ctx context.Context, req *dto.CreateOrUpdateMemberRequest, id string) (*dto.CreateOrUpdateMemberResponse, error)
	GetDetailMember(ctx context.Context, id string) (*dto.GetListMember, error)
	RemoveMember(ctx context.Context, id string) error
}
