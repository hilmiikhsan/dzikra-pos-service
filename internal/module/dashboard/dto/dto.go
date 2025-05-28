package dto

type GetDashboardResponse struct {
	TotalAmount       int `json:"total_amount"`
	TotalExpenses     int `json:"total_expenses"`
	TotalTransactions int `json:"total_transaction"`
	TotalSelling      int `json:"total_selling_product"`
	TotalCapital      int `json:"total_capital"`
	NetSales          int `json:"netsales"`
	ProfitLoss        int `json:"profit_loss"`
}
