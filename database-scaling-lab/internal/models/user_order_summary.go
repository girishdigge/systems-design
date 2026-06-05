package models

import "github.com/google/uuid"

type UserOrderSummary struct {
	UserID      uuid.UUID `json:"user_id"`
	TotalOrders int64     `json:"total_orders"`
}
