package rest

import (
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/response"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *productHandler) createProduct(c *fiber.Ctx) error {
	var (
		req = new(dto.CreateOrUpdateProductRequest)
		ctx = c.Context()
	)

	name := c.FormValue("name")
	req.Name = name
	description := c.FormValue("desc")
	req.Description = description
	realPrice := c.FormValue("real_price")

	realPriceInt, err := strconv.Atoi(realPrice)
	if err != nil {
		log.Warn().Err(err).Msg("handler::createProduct - Invalid real_price")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid real_price"))
	}

	req.RealPrice = realPriceInt
	categoryIDStr := c.FormValue("category_id")

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		log.Warn().Err(err).Msg("handler::createProduct - Invalid category_id")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid category_id"))
	}
	req.CategoryID = categoryID

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::createProduct - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to parse request body"))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::createProduct - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	mf, err := c.MultipartForm()
	if err != nil {
		log.Error().Err(err).Msg("handler::createProduct - Failed to parse multipart form")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid multipart form data"))
	}

	fileHeaders := mf.File[constants.MultipartFormFile]
	switch len(fileHeaders) {
	case 0:
		log.Error().Msg("handler::createProduct - No image file uploaded")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("At least one image file is required"))
	case 1:
		log.Info().Msgf("handler::createProduct - %s file is valid", fileHeaders[0].Filename)
	default:
		log.Error().Msgf("handler::createProduct - too many files uploaded: %d", len(fileHeaders))
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Only one image file is allowed"))
	}

	fh := fileHeaders[0]
	if fh.Size > constants.MaxFileSize {
		log.Error().Msg("handler::createProduct - File size exceeds limit")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("File size exceeds limit"))
	}

	ext := strings.ToLower(filepath.Ext(fh.Filename))
	if !constants.AllowedImageExtensions[ext] {
		log.Error().Msg("handler::createProduct - Invalid file extension")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid file extension"))
	}

	file, err := fh.Open()
	if err != nil {
		log.Error().Err(err).Msg("handler::createProduct - Failed to open file")
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error("Internal server error"))
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Error().Err(err).Msg("handler::createProduct - Failed to read file")
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error("Internal server error"))
	}

	mimeType := http.DetectContentType(fileBytes)
	if !strings.HasPrefix(mimeType, "image/") {
		log.Error().Msg("handler::createProduct - Uploaded file is not a valid image")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Uploaded file is not a valid image"))
	}

	objectName := "product_pos_images/" + utils.GenerateBucketFileUUID() + ext
	uploadFile := dto.UploadFileRequest{
		ObjectName:     objectName,
		File:           fileBytes,
		FileHeaderSize: fh.Size,
		ContentType:    mimeType,
		Filename:       fh.Filename,
	}

	res, err := h.service.CreateProduct(ctx, uploadFile, req)
	if err != nil {
		log.Error().Err(err).Msg("handler::createProduct - Failed to create product content")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}

func (h *productHandler) updateProduct(c *fiber.Ctx) error {
	var (
		req = new(dto.CreateOrUpdateProductRequest)
		ctx = c.Context()
		id  = c.Params("product_id")
	)

	name := c.FormValue("name")
	req.Name = name
	description := c.FormValue("desc")
	req.Description = description
	realPrice := c.FormValue("real_price")

	realPriceInt, err := strconv.Atoi(realPrice)
	if err != nil {
		log.Warn().Err(err).Msg("handler::updateProduct - Invalid real_price")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid real_price"))
	}

	req.RealPrice = realPriceInt
	categoryIDStr := c.FormValue("category_id")

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		log.Warn().Err(err).Msg("handler::updateProduct - Invalid category_id")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid category_id"))
	}
	req.CategoryID = categoryID

	if id == "" {
		log.Warn().Msg("handler::updateProduct - Product ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Product ID is required"))
	}

	productID, err := strconv.Atoi(id)
	if err != nil {
		log.Warn().Err(err).Msg("handler::updateProduct - Invalid product ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid product ID"))
	}

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateProduct - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to parse request body"))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateProduct - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	var uploadFile *dto.UploadFileRequest
	if mf, err := c.MultipartForm(); err == nil {
		files := mf.File[constants.MultipartFormFile]
		if len(files) > 1 {
			log.Error().Msgf("handler::updateProduct - too many files uploaded: %d", len(files))
			return c.Status(fiber.StatusBadRequest).JSON(response.Error("Only one image file is allowed"))
		}
		if len(files) == 1 {
			fh := files[0]
			if fh.Size > constants.MaxFileSize {
				log.Error().Msg("handler::updateProduct - File size exceeds limit")
				return c.Status(fiber.StatusBadRequest).JSON(response.Error("File size exceeds limit"))
			}
			ext := strings.ToLower(filepath.Ext(fh.Filename))
			if !constants.AllowedImageExtensions[ext] {
				log.Error().Msg("handler::updateProduct - Invalid file extension")
				return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid file extension"))
			}
			f, err := fh.Open()
			if err != nil {
				log.Error().Err(err).Msg("handler::updateProduct - Failed to open file")
				return c.Status(fiber.StatusInternalServerError).JSON(response.Error("Internal server error"))
			}
			defer f.Close()
			data, err := io.ReadAll(f)
			if err != nil {
				log.Error().Err(err).Msg("handler::updateProduct - Failed to read file")
				return c.Status(fiber.StatusInternalServerError).JSON(response.Error("Internal server error"))
			}
			mime := http.DetectContentType(data)
			if !strings.HasPrefix(mime, "image/") {
				log.Error().Msg("handler::updateProduct - Uploaded file is not a valid image")
				return c.Status(fiber.StatusBadRequest).JSON(response.Error("Uploaded file is not a valid image"))
			}
			objectName := "product_pos_images/" + utils.GenerateBucketFileUUID() + ext
			uploadFile = &dto.UploadFileRequest{
				ObjectName:     objectName,
				File:           data,
				FileHeaderSize: fh.Size,
				ContentType:    mime,
				Filename:       fh.Filename,
			}
		}
	}

	res, err := h.service.UpdateProduct(ctx, uploadFile, req, productID)
	if err != nil {
		log.Error().Err(err).Msg("handler::updateProduct - Failed to update product")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *productHandler) getListProduct(c *fiber.Ctx) error {
	var (
		ctx    = c.Context()
		page   = c.QueryInt("page", 1)
		limit  = c.QueryInt("limit", 10)
		search = c.Query("search", "")
	)

	res, err := h.service.GetListProduct(ctx, page, limit, search)
	if err != nil {
		log.Error().Err(err).Msg("handler::getListProduct - Failed to get list product")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *productHandler) getDetailProduct(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		id  = c.Params("product_id")
	)

	if id == "" {
		log.Warn().Msg("handler::getDetailProduct - Product ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Product ID is required"))
	}

	productID, err := strconv.Atoi(id)
	if err != nil {
		log.Warn().Err(err).Msg("handler::getDetailProduct - Invalid product ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid product ID"))
	}

	res, err := h.service.GetDetailProduct(ctx, productID)
	if err != nil {
		log.Error().Err(err).Msg("handler::getDetailProduct - Failed to get detail product")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *productHandler) removeProduct(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		id  = c.Params("product_id")
	)

	if id == "" {
		log.Warn().Msg("handler::removeProduct - Product ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Product ID is required"))
	}

	productID, err := strconv.Atoi(id)
	if err != nil {
		log.Warn().Err(err).Msg("handler::removeProduct - Invalid product ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid product ID"))
	}

	err = h.service.RemoveProduct(ctx, productID)
	if err != nil {
		log.Error().Err(err).Msg("handler::removeProduct - Failed to remove product")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success("OK", ""))
}
