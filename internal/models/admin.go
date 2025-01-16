package models

// AdminDashboardStats содержит основные статистические данные
type AdminDashboardStats struct {
	TotalOrders     int64 `json:"total_orders"`
	TotalProducts   int64 `json:"total_products"`
	TotalUsers      int64 `json:"total_users"`
	PendingOrders   int64 `json:"pending_orders"`
	CompletedOrders int64 `json:"completed_orders"`
}

