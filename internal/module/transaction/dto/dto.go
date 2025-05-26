package dto

type CreateTransactionRequest struct {
	Name                     string               `json:"name" validate:"required,min=1,max=20,xss_safe"`
	Status                   string               `json:"status"`
	Email                    string               `json:"email" validate:"required,email,email_blacklist"`
	IsMember                 bool                 `json:"is_member"`
	PhoneNumber              string               `json:"number_phone" validate:"required,phone,max=17,xss_safe"`
	TransactionRequest       []TransactionRequest `json:"transaction_requests"`
	CallbackFinish           string               `json:"callback_finish" validate:"required,callback_finish"`
	TableNumber              string               `json:"table_number" validate:"required,min=1,max=20,xss_safe"`
	Notes                    string               `json:"notes" validate:"required,min=1,max=20,xss_safe"`
	PaymentType              string               `json:"payment_type"`
	TotalMoney               int                  `json:"total_money"`
	TotalQuantity            int                  `json:"total_quantity"`
	TotalProductAmount       int                  `json:"total_product_amount"`
	TotalAmount              int                  `json:"total_amount"`
	DiscountPercentage       int                  `json:"discount_percentage"`
	ChangeMoney              int                  `json:"change_money"`
	TotalProductCapitalPrice int                  `json:"total_product_capital_price"`
	TaxAmount                int                  `json:"tax_amount"`
	TransactionItems         []TransactionItem    `json:"transaction_items"`
}

type CreateTransactionResponse struct {
	ID                       string                    `json:"id"`
	Status                   string                    `json:"status"`
	PhoneNumber              string                    `json:"number_phone"`
	Name                     string                    `json:"name"`
	Email                    string                    `json:"email"`
	IsMember                 bool                      `json:"is_member"`
	TotalQuantity            string                    `json:"total_quantity"`
	TotalProductAmount       string                    `json:"total_product_amount"`
	TotalProductCapitalPrice string                    `json:"total_product_capital_price"`
	TotalAmount              string                    `json:"total_amount"`
	DiscountPercentage       string                    `json:"discount_percentage"`
	VTransactionID           string                    `json:"v_transaction_id"`
	VPaymentID               string                    `json:"v_payment_id"`
	VPaymentRedirectUrl      string                    `json:"v_payment_redirect_url"`
	PaymentType              string                    `json:"payment_type"`
	TotalMoney               *string                   `json:"total_money"`
	ChangeMoney              *string                   `json:"change_money"`
	TableNumber              string                    `json:"table_number"`
	CreatedAt                string                    `json:"created_at"`
	Notes                    string                    `json:"notes"`
	TaxAmount                string                    `json:"tax_amount"`
	TransactionItems         []TransactionItemResponse `json:"transaction_items"`
}

type TransactionRequest struct {
	Quantity  int `json:"quantity"`
	ProductID int `json:"product_id"`
}

type TransactionItem struct {
	ID                      int    `json:"id"`
	Quantity                int    `json:"quantity"`
	TotalAmount             int    `json:"total_amount"`
	ProductName             string `json:"product_name"`
	ProductPrice            int    `json:"product_price"`
	TransactionID           string `json:"transaction_id"`
	ProductID               int    `json:"product_id"`
	TotalAmountCapitalPrice int    `json:"total_amount_capital_price"`
	ProductCapitalPrice     int    `json:"product_capital_price"`
}

type TransactionItemResponse struct {
	ID                      int    `json:"id"`
	Quantity                string `json:"quantity"`
	TotalAmount             string `json:"total_amount"`
	ProductName             string `json:"product_name"`
	ProductPrice            string `json:"product_price"`
	TransactionID           string `json:"transaction_id"`
	ProductID               int    `json:"product_id"`
	TotalAmountCapitalPrice string `json:"total_amount_capital_price"`
	ProductCapitalPrice     string `json:"product_capital_price"`
}

type MinIngredientInput struct {
	Quantity  int `json:"quantity"`
	ProductID int `json:"product_id"`
}

type GetListTransactionResponse struct {
	Transactions []GetListTransaction `json:"transactions"`
	TotalPages   int                  `json:"total_page"`
	CurrentPage  int                  `json:"current_page"`
	PageSize     int                  `json:"page_size"`
	TotalData    int                  `json:"total_data"`
}

type GetListTransaction struct {
	ID                       string `json:"id"`
	Status                   string `json:"status"`
	PhonenUmber              string `json:"number_phone"`
	Name                     string `json:"name"`
	Email                    string `json:"email"`
	IsMember                 bool   `json:"is_member"`
	TotalQuantity            string `json:"total_quantity"`
	TotalProductAmount       string `json:"total_product_amount"`
	TotalProductCapitalPrice string `json:"total_product_capital_price"`
	TotalAmount              string `json:"total_amount"`
	DiscountPercentage       string `json:"discount_percentage"`
	VTransactionID           string `json:"v_transaction_id"`
	VPaymentID               string `json:"v_payment_id"`
	VPaymentRedirectUrl      string `json:"v_payment_redirect_url"`
	PaymentType              string `json:"payment_type"`
	TableNumber              string `json:"table_number"`
	CreatedAt                string `json:"created_at"`
	Notes                    string `json:"notes"`
	TaxAmount                string `json:"tax_amount"`
}

type GetTransactionDetailResponse struct {
	ID                       string                    `json:"id"`
	Status                   string                    `json:"status"`
	PhoneNumber              string                    `json:"number_phone"`
	Name                     string                    `json:"name"`
	Email                    string                    `json:"email"`
	IsMember                 bool                      `json:"is_member"`
	TotalQuantity            string                    `json:"total_quantity"`
	TotalProductAmount       string                    `json:"total_product_amount"`
	TotalProductCapitalPrice string                    `json:"total_product_capital_price"`
	TotalAmount              string                    `json:"total_amount"`
	DiscountPercentage       string                    `json:"discount_percentage"`
	VTransactionID           string                    `json:"v_transaction_id"`
	VPaymentID               string                    `json:"v_payment_id"`
	VPaymentRedirectUrl      string                    `json:"v_payment_redirect_url"`
	PaymentType              string                    `json:"payment_type"`
	TableNumber              string                    `json:"table_number"`
	CreatedAt                string                    `json:"created_at"`
	Notes                    string                    `json:"notes"`
	TaxAmount                string                    `json:"tax_amount"`
	TransactionItem          []TransactionItemResponse `json:"TransactionItem"`
}
