package service

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/infrastructure/config"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product/dto"
	product "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product/entity"
	productCategoryDto "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product_category/dto"
	recipe "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/recipe/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/microcosm-cc/bluemonday"
	"github.com/rs/zerolog/log"
)

func (s *productService) CreateProduct(ctx context.Context, payloadFile dto.UploadFileRequest, req *dto.CreateOrUpdateProductRequest) (*dto.CreateOrUpdateProductResponse, error) {
	// mapping file upload
	ext := strings.ToLower(filepath.Ext(payloadFile.Filename))
	objectName := fmt.Sprintf("product_pos_images/%s_%s", utils.GenerateBucketFileUUID(), ext)
	byteFile := utils.NewByteFile(payloadFile.File)

	productCategoryResult, err := s.productCategoryRepository.FindProductCategoryByID(ctx, req.CategoryID)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrProductCategoryNotFound) {
			log.Error().Err(err).Any("payload", req).Msg("service::CreateProduct - Product category not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrProductCategoryNotFound))
		}

		log.Error().Err(err).Any("payload", req).Msg("service::CreateProduct - Failed to find product category")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateProduct - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Any("payload", req).Msg("service::CreateProduct - Failed to rollback transaction")
			}
		}
	}()

	recipeResult, err := s.recipeRepository.InsertNewRecipe(ctx, tx, &recipe.Recipe{
		CapitalPrice: 0,
	})
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateProduct - Failed to insert new recipe")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	payload := &product.Product{
		Name:              req.Name,
		RealPrice:         req.RealPrice,
		Description:       req.Description,
		Stock:             0,
		ImageUrl:          objectName,
		IsActive:          true,
		ProductCategoryID: productCategoryResult.ID,
		RecipeID:          recipeResult.ID,
	}

	res, err := s.productRepository.InsertNewProduct(ctx, tx, payload)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateProduct - Failed to create new product")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// upload file to minio
	uploadedPath, err := s.minioService.UploadFile(ctx, objectName, byteFile, payloadFile.FileHeaderSize, payloadFile.ContentType)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateProduct - Failed to upload file to minio")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	log.Info().Msgf("Uploaded image URL: %s", uploadedPath)

	// commit transaction
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::CreateProduct - failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	publicURL := config.Envs.MinioStorage.PublicURL
	response := &dto.CreateOrUpdateProductResponse{
		ID:          res.ID,
		Name:        res.Name,
		Description: res.Description,
		Stock:       res.Stock,
		RealPrice:   res.RealPrice,
		ImageUrl:    utils.FormatMediaPathURL(res.ImageUrl, publicURL),
		IsActive:    res.IsActive,
		ProductCategory: productCategoryDto.GetProductCategoryResponse{
			ID:   res.ProductCategoryID,
			Name: productCategoryResult.Name,
		},
	}

	// Sanitize response
	policy := bluemonday.UGCPolicy()
	sanitizedResponse := utils.SanitizeCreateOrUpdateProductResponse(*response, policy)

	return &sanitizedResponse, nil
}

func (s *productService) UpdateProduct(ctx context.Context, payloadFile *dto.UploadFileRequest, req *dto.CreateOrUpdateProductRequest, id int) (*dto.CreateOrUpdateProductResponse, error) {
	countCat, err := s.productRepository.CountProductByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::UpdateProduct - Failed to count product by id")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	if countCat == 0 {
		log.Error().Err(err).Any("payload", req).Msg("service::UpdateProduct - Product not found")
		return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrProductNotFound))
	}

	productCategoryResult, err := s.productCategoryRepository.FindProductCategoryByID(ctx, req.CategoryID)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrProductCategoryNotFound) {
			log.Error().Err(err).Any("payload", req).Msg("service::UpdateProduct - Product category not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrProductCategoryNotFound))
		}

		log.Error().Err(err).Any("payload", req).Msg("service::UpdateProduct - Failed to find product category")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	existing, err := s.productRepository.FindProductByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrProductNotFound) {
			log.Error().Err(err).Msg("service::UpdateProduct - product not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrProductNotFound))
		}

		log.Error().Err(err).Msg("service::UpdateProduct - Failed to get product")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	var objectName string
	if payloadFile != nil {
		ext := strings.ToLower(filepath.Ext(payloadFile.Filename))
		objectName = fmt.Sprintf("product_pos_images/%s%s", utils.GenerateBucketFileUUID(), ext)

		if existing.ImageUrl != "" {
			if err := s.minioService.DeleteFile(ctx, existing.ImageUrl); err != nil {
				log.Error().Err(err).Msg("service::UpdateProduct - Failed to delete old product image")
				return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
			}
		}
	}

	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::UpdateProduct - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Any("payload", req).Msg("service::UpdateProduct - Failed to rollback transaction")
			}
		}
	}()

	payload := &product.Product{
		ID:                id,
		Name:              req.Name,
		RealPrice:         req.RealPrice,
		Description:       req.Description,
		Stock:             0,
		ImageUrl:          existing.ImageUrl,
		IsActive:          true,
		ProductCategoryID: productCategoryResult.ID,
		RecipeID:          existing.RecipeID,
	}

	if objectName != "" {
		payload.ImageUrl = objectName
	}

	res, err := s.productRepository.UpdateProduct(ctx, tx, payload)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrProductNotFound) {
			log.Error().Err(err).Msg("service::UpdateProduct - article not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrProductNotFound))
		}

		log.Error().Err(err).Any("payload", req).Msg("service::UpdateProduct - Failed to update product")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	if payloadFile != nil {
		byteFile := utils.NewByteFile(payloadFile.File)
		if _, err = s.minioService.UploadFile(ctx, objectName, byteFile, payloadFile.FileHeaderSize, payloadFile.ContentType); err != nil {
			log.Error().Err(err).Msg("service::UpdateProduct - Failed to upload file to minio")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}
	}

	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::UpdateProduct - failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	publicURL := config.Envs.MinioStorage.PublicURL
	response := &dto.CreateOrUpdateProductResponse{
		ID:          res.ID,
		Name:        res.Name,
		Description: res.Description,
		Stock:       res.Stock,
		RealPrice:   res.RealPrice,
		ImageUrl:    utils.FormatMediaPathURL(res.ImageUrl, publicURL),
		IsActive:    res.IsActive,
		ProductCategory: productCategoryDto.GetProductCategoryResponse{
			ID:   res.ProductCategoryID,
			Name: productCategoryResult.Name,
		},
	}

	// Sanitize response
	policy := bluemonday.UGCPolicy()
	sanitizedResponse := utils.SanitizeCreateOrUpdateProductResponse(*response, policy)

	return &sanitizedResponse, nil
}

func (s *productService) GetListProduct(ctx context.Context, page, limit int, search string) (*dto.GetListProductResponse, error) {
	// calculate pagination
	currentPage, perPage, offset := utils.Paginate(page, limit)

	// get list products
	productss, total, err := s.productRepository.FindListProduct(ctx, perPage, offset, search)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListProduct - error getting list product")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check if productss is nil
	if productss == nil {
		productss = []dto.GetListProduct{}
	}

	// calculate total pages
	totalPages := utils.CalculateTotalPages(total, perPage)

	// create map response
	response := dto.GetListProductResponse{
		Product:     productss,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
		PageSize:    perPage,
		TotalData:   total,
	}

	// return response
	return &response, nil
}

func (s *productService) GetDetailProduct(ctx context.Context, id int) (*dto.GetListProduct, error) {
	res, err := s.productRepository.FindProductByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrProductNotFound) {
			log.Error().Err(err).Msg("service::GetDetailProduct - product not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrProductNotFound))
		}

		log.Error().Err(err).Msg("service::GetDetailProduct - Failed to get product")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	publicURL := config.Envs.MinioStorage.PublicURL
	response := &dto.GetListProduct{
		ID:          res.ID,
		Name:        res.Name,
		Description: res.Description,
		Stock:       res.Stock,
		RealPrice:   res.RealPrice,
		ImageUrl:    utils.FormatMediaPathURL(res.ImageUrl, publicURL),
		IsActive:    res.IsActive,
		ProductCategory: productCategoryDto.GetProductCategoryResponse{
			ID:   res.ProductCategoryID,
			Name: res.ProductCategoryName,
		},
	}

	return response, nil
}

func (s *productService) RemoveProduct(ctx context.Context, id int) error {
	product, err := s.productRepository.FindProductByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrProductNotFound) {
			log.Error().Err(err).Msg("service::RemoveProduct - product not found")
			return err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrProductNotFound))
		}

		log.Error().Err(err).Msg("service::RemoveProduct - error getting product")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Msg("service::RemoveProduct - Failed to begin transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Msg("service::RemoveProduct - Failed to rollback transaction")
			}
		}
	}()

	if err := s.productRepository.SoftDeleteProductByID(ctx, tx, id); err != nil {
		log.Error().Err(err).Msg("service::RemoveProduct - error removing product")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	if product.ImageUrl != "" {
		if err := s.minioService.DeleteFile(ctx, product.ImageUrl); err != nil {
			log.Error().Err(err).Msg("service::RemoveProduct - Failed to delete product image")
			return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::RemoveProduct - failed to commit transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return nil
}
