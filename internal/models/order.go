package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID            uuid.UUID   `json:"id"`
	CreatorID     uuid.UUID   `json:"creator_id"`
	CustomerName  string      `json:"customer_name"`
	CustomerPhone string      `json:"customer_phone"`
	Items         []OrderItem `json:"items"`
	TotalAmount   float64     `json:"total_amount"`
	Status        string      `json:"status"` // "pending", "paid", "cancelled"
	StatusAr      string      `json:"status_ar"`
	OrderSource   string      `json:"order_source"` // "whatsapp", "direct"
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID             uuid.UUID `json:"id"`
	OrderID        uuid.UUID `json:"order_id"`
	WorkshopID     uuid.UUID `json:"workshop_id"`
	WorkshopName   string    `json:"workshop_name"`
	WorkshopNameAr string    `json:"workshop_name_ar"`
	Price          float64   `json:"price"`
	Quantity       int       `json:"quantity"`
}

type OrderStats struct {
	PendingOrders   int     `json:"pending_orders"`
	PaidOrders      int     `json:"paid_orders"`
	CancelledOrders int     `json:"cancelled_orders"`
	TotalRevenue    float64 `json:"total_revenue"`
}

type CreateOrderRequest struct {
	CustomerName    string             `json:"customer_name"`
	CustomerPhone   string             `json:"customer_phone"`
	Items           []OrderItemRequest `json:"items"`
	OrderSource     string             `json:"order_source"`
	CreatorUsername string             `json:"creator_username"`
}

type OrderItemRequest struct {
	WorkshopID uuid.UUID `json:"workshop_id"`
	Quantity   int       `json:"quantity"`
}
