package dto

type CreateOrUpdateExpensesRequest struct {
	Name string `json:"name" validate:"required,min=1,max=30,xss_safe"`
	Cost int    `json:"cost" validate:"required,numeric,number"`
	Date string `json:"created_at" validate:"required,date_format"`
}

type CreateOrUpdateExpensesResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Cost int    `json:"cost"`
	Date string `json:"created_at"`
}

type GetListExpensesResponse struct {
	Expenses    []GetListExpenses `json:"expenses"`
	TotalPages  int               `json:"total_page"`
	CurrentPage int               `json:"current_page"`
	PageSize    int               `json:"page_size"`
	TotalData   int               `json:"total_data"`
}

type GetListExpenses struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Cost      int    `json:"cost"`
	Date      string `json:"date"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
