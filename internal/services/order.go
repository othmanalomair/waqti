package services

import (
	"fmt"
	"time"
	"waqti/internal/models"

	"github.com/google/uuid"
)

type OrderService struct {
	orders []models.Order
}

func NewOrderService() *OrderService {
	// Generate fixed UUIDs for demo data consistency
	creatorID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	workshop1ID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	workshop2ID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440002")

	orders := []models.Order{
		{
			ID:            uuid.MustParse("550e8400-e29b-41d4-a716-446655440030"),
			CreatorID:     creatorID,
			CustomerName:  "أحمد محمد",
			CustomerPhone: "+965-9999-1234",
			Items: []models.OrderItem{
				{
					ID:             uuid.MustParse("550e8400-e29b-41d4-a716-446655440040"),
					OrderID:        uuid.MustParse("550e8400-e29b-41d4-a716-446655440030"),
					WorkshopID:     workshop1ID,
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
			ID:            uuid.MustParse("550e8400-e29b-41d4-a716-446655440031"),
			CreatorID:     creatorID,
			CustomerName:  "سارة أحمد",
			CustomerPhone: "+965-9999-5678",
			Items: []models.OrderItem{
				{
					ID:             uuid.MustParse("550e8400-e29b-41d4-a716-446655440041"),
					OrderID:        uuid.MustParse("550e8400-e29b-41d4-a716-446655440031"),
					WorkshopID:     workshop2ID,
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

func (s *OrderService) CreateOrder(creatorID uuid.UUID, request models.CreateOrderRequest) (*models.Order, error) {
	newID := uuid.New()

	var totalAmount float64
	var items []models.OrderItem

	workshopService := NewWorkshopService()

	for _, itemReq := range request.Items {
		workshop := s.getWorkshopByID(itemReq.WorkshopID, workshopService)
		if workshop == nil {
			return nil, fmt.Errorf("workshop not found: %s", itemReq.WorkshopID.String())
		}

		item := models.OrderItem{
			ID:             uuid.New(),
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

func (s *OrderService) getWorkshopByID(workshopID uuid.UUID, workshopService *WorkshopService) *models.Workshop {
	workshops := workshopService.GetWorkshopsByCreatorID(uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"))
	for _, workshop := range workshops {
		if workshop.ID == workshopID {
			return &workshop
		}
	}
	return nil
}

func (s *OrderService) GetOrdersByCreatorID(creatorID uuid.UUID, filter models.EnrollmentFilter) []models.Order {
	var filtered []models.Order

	for _, order := range s.orders {
		if order.CreatorID == creatorID {
			filtered = append(filtered, order)
		}
	}

	filteredByTime := s.filterByTimeRange(filtered, filter.TimeRange)
	s.sortOrders(filteredByTime, filter.OrderBy, filter.OrderDir)

	return filteredByTime
}

func (s *OrderService) GetOrderStats(creatorID uuid.UUID, timeRange string) models.OrderStats {
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

func (s *OrderService) GetPendingOrdersCount(creatorID uuid.UUID) int {
	count := 0
	for _, order := range s.orders {
		if order.CreatorID == creatorID && order.Status == "pending" {
			count++
		}
	}
	return count
}

func (s *OrderService) UpdateOrderStatus(orderID uuid.UUID, newStatus string) error {
	for i, order := range s.orders {
		if order.ID == orderID {
			s.orders[i].Status = newStatus
			s.orders[i].UpdatedAt = time.Now()

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

func (s *OrderService) DeleteOrder(orderID uuid.UUID) error {
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
		cutoff = now.AddDate(0, 0, -30)
	case "months":
		cutoff = now.AddDate(0, -12, 0)
	case "year":
		cutoff = now.AddDate(-5, 0, 0)
	default:
		cutoff = now.AddDate(0, 0, -30)
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
	// Implementation for sorting orders
}
