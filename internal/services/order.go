package services

import (
	"fmt"
	"time"
	"waqti/internal/database"
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

	// Validate all workshops exist and calculate total
	for _, itemReq := range request.Items {
		workshop, err := workshopService.GetWorkshopByID(itemReq.WorkshopID, creatorID)
		if err != nil || workshop == nil {
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

	// Try to save to database first
	err := s.saveOrderToDatabase(&order)
	if err != nil {
		// If database fails, fall back to in-memory storage for demo
		fmt.Printf("Database order creation failed, using fallback: %v\n", err)
		s.orders = append(s.orders, order)
	}

	return &order, nil
}

func (s *OrderService) saveOrderToDatabase(order *models.Order) error {
	// Start transaction
	tx, err := database.Instance.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert order
	orderQuery := `
		INSERT INTO orders (
			id, creator_id, customer_name, customer_phone, 
			total_amount, status, order_source, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err = tx.Exec(orderQuery, 
		order.ID, order.CreatorID, order.CustomerName, order.CustomerPhone,
		order.TotalAmount, order.Status, order.OrderSource, 
		order.CreatedAt, order.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	// Insert order items
	itemQuery := `
		INSERT INTO order_items (
			id, order_id, workshop_id, workshop_name, workshop_name_ar,
			price, quantity, subtotal
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	for _, item := range order.Items {
		subtotal := item.Price * float64(item.Quantity)
		_, err = tx.Exec(itemQuery,
			item.ID, item.OrderID, item.WorkshopID, item.WorkshopName, item.WorkshopNameAr,
			item.Price, item.Quantity, subtotal,
		)
		if err != nil {
			return fmt.Errorf("failed to insert order item: %w", err)
		}
	}

	// Commit transaction
	return tx.Commit()
}


func (s *OrderService) GetOrdersByCreatorID(creatorID uuid.UUID, filter models.EnrollmentFilter) []models.Order {
	// Try to get orders from database first
	dbOrders, err := s.getOrdersFromDatabase(creatorID, filter)
	if err != nil {
		fmt.Printf("Database order retrieval failed, using fallback: %v\n", err)
		// Fall back to in-memory data
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

	return dbOrders
}

func (s *OrderService) getOrdersFromDatabase(creatorID uuid.UUID, filter models.EnrollmentFilter) ([]models.Order, error) {
	// Build the base query
	query := `
		SELECT o.id, o.creator_id, o.customer_name, o.customer_phone, 
		       o.total_amount, o.status, o.order_source, o.created_at, o.updated_at
		FROM orders o
		WHERE o.creator_id = $1
	`

	// Add time filter
	args := []interface{}{creatorID}
	argIndex := 2

	now := time.Now()
	var cutoff time.Time
	switch filter.TimeRange {
	case "days":
		cutoff = now.AddDate(0, 0, -30)
	case "months":
		cutoff = now.AddDate(0, -12, 0)
	case "year":
		cutoff = now.AddDate(-5, 0, 0)
	default:
		cutoff = now.AddDate(0, 0, -30)
	}

	query += fmt.Sprintf(" AND o.created_at >= $%d", argIndex)
	args = append(args, cutoff)
	argIndex++

	// Add ordering
	orderBy := "o.created_at"
	switch filter.OrderBy {
	case "price":
		orderBy = "o.total_amount"
	case "name":
		orderBy = "o.customer_name"
	default:
		orderBy = "o.created_at"
	}

	orderDir := "DESC"
	if filter.OrderDir == "asc" {
		orderDir = "ASC"
	}

	query += fmt.Sprintf(" ORDER BY %s %s", orderBy, orderDir)

	rows, err := database.Instance.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID, &order.CreatorID, &order.CustomerName, &order.CustomerPhone,
			&order.TotalAmount, &order.Status, &order.OrderSource,
			&order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		// Set Arabic status
		switch order.Status {
		case "pending":
			order.StatusAr = "قيد الانتظار"
		case "paid":
			order.StatusAr = "مدفوع"
		case "cancelled":
			order.StatusAr = "ملغي"
		}

		// Get order items
		items, err := s.getOrderItemsFromDatabase(order.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get order items: %w", err)
		}
		order.Items = items

		orders = append(orders, order)
	}

	return orders, nil
}

func (s *OrderService) getOrderItemsFromDatabase(orderID uuid.UUID) ([]models.OrderItem, error) {
	query := `
		SELECT id, order_id, workshop_id, workshop_name, workshop_name_ar, price, quantity
		FROM order_items
		WHERE order_id = $1
	`

	rows, err := database.Instance.Query(query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to query order items: %w", err)
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		err := rows.Scan(
			&item.ID, &item.OrderID, &item.WorkshopID, &item.WorkshopName, &item.WorkshopNameAr,
			&item.Price, &item.Quantity,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

func (s *OrderService) GetOrderStats(creatorID uuid.UUID, timeRange string) models.OrderStats {
	// Try to get stats from database first
	stats, err := s.getOrderStatsFromDatabase(creatorID, timeRange)
	if err != nil {
		fmt.Printf("Database order stats failed, using fallback: %v\n", err)
		// Fall back to in-memory data
		orders := s.filterByTimeRange(s.orders, timeRange)
		fallbackStats := models.OrderStats{}

		for _, order := range orders {
			if order.CreatorID != creatorID {
				continue
			}

			switch order.Status {
			case "pending":
				fallbackStats.PendingOrders++
			case "paid":
				fallbackStats.PaidOrders++
				fallbackStats.TotalRevenue += order.TotalAmount
			case "cancelled":
				fallbackStats.CancelledOrders++
			}
		}
		return fallbackStats
	}

	return stats
}

func (s *OrderService) getOrderStatsFromDatabase(creatorID uuid.UUID, timeRange string) (models.OrderStats, error) {
	// Calculate time cutoff
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

	query := `
		SELECT 
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_orders,
			COUNT(CASE WHEN status = 'paid' THEN 1 END) as paid_orders,
			COUNT(CASE WHEN status = 'cancelled' THEN 1 END) as cancelled_orders,
			COALESCE(SUM(CASE WHEN status = 'paid' THEN total_amount ELSE 0 END), 0) as total_revenue
		FROM orders 
		WHERE creator_id = $1 AND created_at >= $2
	`
	
	var stats models.OrderStats
	err := database.Instance.QueryRow(query, creatorID, cutoff).Scan(
		&stats.PendingOrders,
		&stats.PaidOrders,
		&stats.CancelledOrders,
		&stats.TotalRevenue,
	)
	if err != nil {
		return models.OrderStats{}, fmt.Errorf("failed to get order stats: %w", err)
	}
	
	return stats, nil
}

func (s *OrderService) GetPendingOrdersCount(creatorID uuid.UUID) int {
	// Try to get count from database first
	count, err := s.getPendingOrdersCountFromDatabase(creatorID)
	if err != nil {
		fmt.Printf("Database pending orders count failed, using fallback: %v\n", err)
		// Fall back to in-memory data
		fallbackCount := 0
		for _, order := range s.orders {
			if order.CreatorID == creatorID && order.Status == "pending" {
				fallbackCount++
			}
		}
		return fallbackCount
	}
	return count
}

func (s *OrderService) getPendingOrdersCountFromDatabase(creatorID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM orders 
		WHERE creator_id = $1 AND status = 'pending'
	`
	
	var count int
	err := database.Instance.QueryRow(query, creatorID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count pending orders: %w", err)
	}
	
	return count, nil
}

func (s *OrderService) UpdateOrderStatus(orderID uuid.UUID, newStatus string) error {
	// Try to update in database first
	err := s.updateOrderStatusInDatabase(orderID, newStatus)
	if err != nil {
		fmt.Printf("Database order status update failed, using fallback: %v\n", err)
		// Fall back to in-memory update
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
	return nil
}

func (s *OrderService) updateOrderStatusInDatabase(orderID uuid.UUID, newStatus string) error {
	query := `
		UPDATE orders 
		SET status = $1, updated_at = $2 
		WHERE id = $3
	`
	
	result, err := database.Instance.Exec(query, newStatus, time.Now(), orderID)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("order not found")
	}
	
	return nil
}

func (s *OrderService) DeleteOrder(orderID uuid.UUID) error {
	// Try to delete from database first
	err := s.deleteOrderFromDatabase(orderID)
	if err != nil {
		fmt.Printf("Database order deletion failed, using fallback: %v\n", err)
		// Fall back to in-memory deletion
		for i, order := range s.orders {
			if order.ID == orderID {
				s.orders = append(s.orders[:i], s.orders[i+1:]...)
				return nil
			}
		}
		return fmt.Errorf("order not found")
	}
	return nil
}

func (s *OrderService) deleteOrderFromDatabase(orderID uuid.UUID) error {
	// Start transaction to delete order and its items
	tx, err := database.Instance.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete order items first (due to foreign key constraint)
	_, err = tx.Exec("DELETE FROM order_items WHERE order_id = $1", orderID)
	if err != nil {
		return fmt.Errorf("failed to delete order items: %w", err)
	}

	// Delete the order
	result, err := tx.Exec("DELETE FROM orders WHERE id = $1", orderID)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("order not found")
	}

	// Commit transaction
	return tx.Commit()
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
