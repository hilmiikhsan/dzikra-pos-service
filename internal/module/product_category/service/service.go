package service

import (
	"context"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product_category/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product_category/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (s *productCategoryService) CreateOrUpdateProductCategory(ctx context.Context, req *dto.CreateOrUpdateProductCategoryRequest) (*dto.CreateOrUpdateProductCategoryResponse, error) {
	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateOrUpdateProductCategory - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Any("payload", req).Msg("service::CreateOrUpdateProductCategory - Failed to rollback transaction")
			}
		}
	}()

	payload := &entity.ProductCategory{
		Name: req.Name,
	}

	res, err := s.productCategoryRepository.InsertNewProductCategory(ctx, tx, payload)
	if err != nil {
		if err.Error() == constants.ErrProductCategoryAlreadyRegistered {
			log.Error().Err(err).Any("payload", req).Msg("service::CreateOrUpdateProductCategory - Product category already exists")
			return nil, err_msg.NewCustomErrors(fiber.StatusConflict, err_msg.WithMessage(constants.ErrProductCategoryAlreadyRegistered))
		}

		log.Error().Err(err).Any("payload", req).Msg("service::CreateOrUpdateProductCategory - Failed to create or update product category")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::CreateOrUpdateProductCategory - failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return &dto.CreateOrUpdateProductCategoryResponse{
		ID:   res.ID,
		Name: res.Name,
	}, nil
}

func (s *productCategoryService) UpdateProductCategory(ctx context.Context, req *dto.CreateOrUpdateProductCategoryRequest, id int) (*dto.CreateOrUpdateProductCategoryResponse, error) {
	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::UpdateProductCategory - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Any("payload", req).Msg("service::UpdateProductCategory - Failed to rollback transaction")
			}
		}
	}()

	payload := &entity.ProductCategory{
		ID:   id,
		Name: req.Name,
	}

	res, err := s.productCategoryRepository.UpdateProductCategory(ctx, tx, payload)
	if err != nil {
		if err.Error() == constants.ErrProductCategoryAlreadyRegistered {
			log.Error().Err(err).Any("payload", req).Msg("service::UpdateProductCategory - Product category already exists")
			return nil, err_msg.NewCustomErrors(fiber.StatusConflict, err_msg.WithMessage(constants.ErrProductCategoryAlreadyRegistered))
		}

		if strings.Contains(err.Error(), constants.ErrProductCategoryNotFound) {
			log.Error().Err(err).Any("payload", req).Msg("service::UpdateProductCategory - Product category not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrProductCategoryNotFound))
		}

		log.Error().Err(err).Any("payload", req).Msg("service::UpdateProductCategory - Failed to create or update product category")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::UpdateProductCategory - failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return &dto.CreateOrUpdateProductCategoryResponse{
		ID:   res.ID,
		Name: res.Name,
	}, nil
}

func (r *productCategoryService) GetDetailProductCategory(ctx context.Context, id int) (*dto.GetProductCategoryResponse, error) {
	res, err := r.productCategoryRepository.FindProductCategoryByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrProductCategoryNotFound) {
			log.Error().Err(err).Msg("service::GetDetailProductCategory - Product category not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrProductCategoryNotFound))
		}

		log.Error().Err(err).Msg("service::GetDetailProductCategory - Failed to get product category")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return &dto.GetProductCategoryResponse{
		ID:   res.ID,
		Name: res.Name,
	}, nil
}

func (s *productCategoryService) GetListProductCategory(ctx context.Context, page, limit int, search string) (*dto.GetListProductCategoryResponse, error) {
	// calculate pagination
	currentPage, perPage, offset := utils.Paginate(page, limit)

	// get list product category
	productCategories, total, err := s.productCategoryRepository.FindListProductCategory(ctx, perPage, offset, search)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListProductCategory - error getting list product category")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check if productCategories is nil
	if productCategories == nil {
		productCategories = []dto.GetProductCategoryResponse{}
	}

	// calculate total pages
	totalPages := utils.CalculateTotalPages(total, perPage)

	// create map response
	response := dto.GetListProductCategoryResponse{
		Category:    productCategories,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
		PageSize:    perPage,
		TotalData:   total,
	}

	// return response
	return &response, nil
}

func (s *productCategoryService) RemoveProductCategory(ctx context.Context, id int) error {
	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Msg("service::RemoveProductCategory - Failed to begin transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Msg("service::RemoveProductCategory - Failed to rollback transaction")
			}
		}
	}()

	err = s.productCategoryRepository.SoftDeleteProductCategoryByID(ctx, tx, id)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrProductCategoryNotFound) {
			log.Error().Err(err).Msg("service::RemoveProductCategory - Product category not found")
			return err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrProductCategoryNotFound))
		}

		log.Error().Err(err).Msg("service::RemoveProductCategory - Failed to remove product category")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::RemoveProductCategory - failed to commit transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return nil
}
