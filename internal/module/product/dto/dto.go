package dto

import (
	product_category "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product_category/dto"
)

type CreateOrUpdateProductRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=30,xss_safe"`
	Description string `json:"desc" validate:"required,min=2,max=100,xss_safe"`
	RealPrice   int    `json:"real_price" validate:"required,numeric,number,gt=0"`
	CategoryID  int    `json:"category_id" validate:"required,numeric,number,gt=0"`
}

type CreateOrUpdateProductResponse struct {
	ID              int                                         `json:"id"`
	Name            string                                      `json:"name"`
	Description     string                                      `json:"desc"`
	Stock           int                                         `json:"stock"`
	RealPrice       int                                         `json:"real_price"`
	CreatedAt       string                                      `json:"created_at"`
	ImageUrl        string                                      `json:"image_url"`
	IsActive        bool                                        `json:"is_active"`
	ProductCategory product_category.GetProductCategoryResponse `json:"product_category"`
}

type UploadFileRequest struct {
	ObjectName     string `json:"object_name"`
	File           []byte `json:"-"`
	FileHeaderSize int64  `json:"-"`
	ContentType    string `json:"-"`
	Filename       string `json:"-"`
}

type GetListProductResponse struct {
	Product     []GetListProduct `json:"product"`
	TotalPages  int              `json:"total_page"`
	CurrentPage int              `json:"current_page"`
	PageSize    int              `json:"page_size"`
	TotalData   int              `json:"total_data"`
}

type GetListProduct struct {
	ID              int                                         `json:"id"`
	Name            string                                      `json:"name"`
	Description     string                                      `json:"desc"`
	Stock           int                                         `json:"stock"`
	RealPrice       int                                         `json:"real_price"`
	CreatedAt       string                                      `json:"created_at"`
	ImageUrl        string                                      `json:"image_url"`
	IsActive        bool                                        `json:"is_active"`
	ProductCategory product_category.GetProductCategoryResponse `json:"product_category"`
}
