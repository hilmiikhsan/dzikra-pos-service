package dto

type CreateOrUpdateProductCategoryRequest struct {
	Name string `json:"category" validate:"required,min=2,max=30,xss_safe"`
}

type CreateOrUpdateProductCategoryResponse struct {
	ID   int    `json:"id"`
	Name string `json:"category"`
}

type GetListProductCategoryResponse struct {
	Category    []GetProductCategoryResponse `json:"category"`
	TotalPages  int                          `json:"total_page"`
	CurrentPage int                          `json:"current_page"`
	PageSize    int                          `json:"page_size"`
	TotalData   int                          `json:"total_data"`
}

type GetProductCategoryResponse struct {
	ID   int    `json:"id"`
	Name string `json:"category"`
}
