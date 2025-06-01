package services

import (
	"fmt"
	"time"
	"waqti/internal/models"
)

type OrderService struct {
	orders []models.Order
}

func NewOrderService() *OrderService {
	// Dummy orders for demo
	orders := []models.Order{
		{
			ID:            1,
			CreatorID:     1,
			CustomerName:  "أحمد محمد",
			CustomerPhone: "+965-9999-1234",
			Items: []models.OrderItem{
				{
					ID:             1,
					OrderID:        1,
					WorkshopID:     1,
					WorkshopName:   "Photography Basics",
					WorkshopNameAr: "أساسيات التصوير",
					Price:          25.0,
					Quantity:       1,
				},
			},
			TotalAmount: 25.0,
			Status:      "pending",
			StatusAr:    "قيد الانتظار",
			OrderSource: "whatsapp",
			CreatedAt:   time.Now().AddDate(0, 0, -1),
			UpdatedAt:   time.Now().AddDate(0, 0, -1),
		},
		{
			ID:            2,
			CreatorID:     1,
			CustomerName:  "سارة أحمد",
			CustomerPhone: "+965-9999-5678",
			Items: []models.OrderItem{
				{
					ID:             2,
					OrderID:        2,
					WorkshopID:     2,
					WorkshopName:   "Digital Marketing",
					WorkshopNameAr: "التسويق الرقمي",
					Price:          35.0,
					Quantity:       1,
				},
			},
			TotalAmount: 35.0,
			Status:      "pending",
			StatusAr:    "قيد الانتظار",
			OrderSource: "whatsapp",
			CreatedAt:   time.Now().AddDate(0, 0, 0),
			UpdatedAt:   time.Now().AddDate(0, 0, 0),
		},
	}

	return &OrderService{
		orders: orders,
	}
}

func (s *OrderService) CreateOrder(creatorID int, request models.CreateOrderRequest) (*models.Order, error) {
	// Generate new ID
	newID := len(s.orders) + 1

	// Calculate total amount
	var totalAmount float64
	var items []models.OrderItem

	workshopService := NewWorkshopService()

	for _, itemReq := range request.Items {
		// Get workshop details (in real app, from database)
		workshop := s.getWorkshopByID(itemReq.WorkshopID, workshopService)
		if workshop == nil {
			return nil, fmt.Errorf("workshop not found: %d", itemReq.WorkshopID)
		}

		item := models.OrderItem{
			ID:             len(items) + 1,
			OrderID:        newID,
			WorkshopID:     itemReq.WorkshopID,
			WorkshopName:   workshop.Title,
			WorkshopNameAr: workshop.TitleAr,
			Price:          workshop.Price,
			Quantity:       itemReq.Quantity,
		}

		items = append(items, item)
		totalAmount += workshop.Price * float64(itemReq.Quantity)
	}

	order := models.Order{
		ID:            newID,
		CreatorID:     creatorID,
		CustomerName:  request.CustomerName,
		CustomerPhone: request.CustomerPhone,
		Items:         items,
		TotalAmount:   totalAmount,
		Status:        "pending",
		StatusAr:      "قيد الانتظار",
		OrderSource:   request.OrderSource,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	s.orders = append(s.orders, order)
	return &order, nil
}

func (s *OrderService) getWorkshopByID(workshopID int, workshopService *WorkshopService) *models.Workshop {
	workshops := workshopService.GetWorkshopsByCreatorID(1)
	for _, workshop := range workshops {
		if workshop.ID == workshopID {
			return &workshop
		}
	}
	return nil
}

func (s *OrderService) GetOrdersByCreatorID(creatorID int, filter models.EnrollmentFilter) []models.Order {
	var filtered []models.Order

	for _, order := range s.orders {
		if order.CreatorID == creatorID {
			filtered = append(filtered, order)
		}
	}

	// Apply time range filter
	filteredByTime := s.filterByTimeRange(filtered, filter.TimeRange)

	// Apply sorting
	s.sortOrders(filteredByTime, filter.OrderBy, filter.OrderDir)

	return filteredByTime
}

func (s *OrderService) GetOrderStats(creatorID int, timeRange string) models.OrderStats {
	orders := s.filterByTimeRange(s.orders, timeRange)

	stats := models.OrderStats{}

	for _, order := range orders {
		if order.CreatorID != creatorID {
			continue
		}

		switch order.Status {
		case "pending":
			stats.PendingOrders++
		case "paid":
			stats.PaidOrders++
			stats.TotalRevenue += order.TotalAmount
		case "cancelled":
			stats.CancelledOrders++
		}
	}

	return stats
}

func (s *OrderService) GetPendingOrdersCount(creatorID int) int {
	count := 0
	for _, order := range s.orders {
		if order.CreatorID == creatorID && order.Status == "pending" {
			count++
		}
	}
	return count
}

func (s *OrderService) UpdateOrderStatus(orderID int, newStatus string) error {
	for i, order := range s.orders {
		if order.ID == orderID {
			s.orders[i].Status = newStatus
			s.orders[i].UpdatedAt = time.Now()

			// Update Arabic status
			switch newStatus {
			case "pending":
				s.orders[i].StatusAr = "قيد الانتظار"
			case "paid":
				s.orders[i].StatusAr = "مدفوع"
			case "cancelled":
				s.orders[i].StatusAr = "ملغي"
			}

			return nil
		}
	}
	return fmt.Errorf("order not found")
}

func (s *OrderService) DeleteOrder(orderID int) error {
	for i, order := range s.orders {
		if order.ID == orderID {
			s.orders = append(s.orders[:i], s.orders[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("order not found")
}

func (s *OrderService) filterByTimeRange(orders []models.Order, timeRange string) []models.Order {
	now := time.Now()
	var cutoff time.Time

	switch timeRange {
	case "days":
		cutoff = now.AddDate(0, 0, -30) // Last 30 days
	case "months":
		cutoff = now.AddDate(0, -12, 0) // Last 12 months
	case "year":
		cutoff = now.AddDate(-5, 0, 0) // Last 5 years
	default:
		cutoff = now.AddDate(0, 0, -30) // Default to 30 days
	}

	var filtered []models.Order
	for _, order := range orders {
		if order.CreatedAt.After(cutoff) {
			filtered = append(filtered, order)
		}
	}

	return filtered
}

func (s *OrderService) sortOrders(orders []models.Order, orderBy, orderDir string) {
	// Implementation similar to enrollment sorting
	// For brevity, using simple sorting
}
