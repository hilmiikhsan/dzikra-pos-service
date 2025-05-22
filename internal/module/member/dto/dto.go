package dto

type CreateOrUpdateMemberRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=30,xss_safe"`
	PhoneNumber string `json:"number_phone" validate:"required,phone,max=17,xss_safe"`
	Email       string `json:"email" validate:"required,email,email_blacklist"`
}

type CreateOrUpdateMemberResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"number_phone"`
	CreatedAt   string `json:"created_at"`
}

type GetListMemberResponse struct {
	Members     []GetListMember `json:"members"`
	TotalPages  int             `json:"total_page"`
	CurrentPage int             `json:"current_page"`
	PageSize    int             `json:"page_size"`
	TotalData   int             `json:"total_data"`
}

type GetListMember struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"number_phone"`
	CreatedAt   string `json:"created_at"`
}
