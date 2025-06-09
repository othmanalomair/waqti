package services

import (
	"database/sql"
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
	sessionService := NewWorkshopSessionService()

	// Validate all workshops exist and calculate total
	for _, itemReq := range request.Items {
		workshop, err := workshopService.GetWorkshopByID(itemReq.WorkshopID, creatorID)
		if err != nil || workshop == nil {
			return nil, fmt.Errorf("workshop not found: %s", itemReq.WorkshopID.String())
		}

		// Find available session for this workshop
		var sessionID *uuid.UUID
		var runID *uuid.UUID
		
		if itemReq.SessionID != nil {
			// Use specific session if provided - validate it exists and get run_id
			session, err := s.getSessionByID(*itemReq.SessionID)
			if err != nil {
				return nil, fmt.Errorf("specified session not found or not available: %s", itemReq.SessionID.String())
			}
			sessionID = itemReq.SessionID
			runID = session.RunID
		} else {
			// Find next available session
			session, err := sessionService.GetNextAvailableSession(itemReq.WorkshopID)
			if err != nil {
				return nil, fmt.Errorf("no available sessions for workshop %s: %w", itemReq.WorkshopID.String(), err)
			}
			sessionID = &session.ID
			runID = session.RunID
		}

		subtotal := workshop.Price * float64(itemReq.Quantity)

		item := models.OrderItem{
			ID:             uuid.New(),
			OrderID:        newID,
			WorkshopID:     itemReq.WorkshopID,
			SessionID:      sessionID,
			RunID:          runID,
			WorkshopName:   workshop.Title,
			WorkshopNameAr: workshop.TitleAr,
			Price:          workshop.Price,
			Quantity:       itemReq.Quantity,
			Subtotal:       subtotal,
		}

		items = append(items, item)
		totalAmount += subtotal
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
	} else {
		// If order saved successfully, create enrollments but DO NOT increment attendance yet
		// Attendance should only be incremented when order is marked as PAID
		enrollmentService := NewEnrollmentService()
		
		for _, item := range order.Items {
			if item.SessionID != nil {
				// Create enrollment records with pending status
				for i := 0; i < item.Quantity; i++ {
					enrollment := &models.Enrollment{
						ID:             uuid.New(),
						WorkshopID:     item.WorkshopID,
						SessionID:      item.SessionID,
						OrderID:        &order.ID,
						WorkshopName:   item.WorkshopName,
						WorkshopNameAr: item.WorkshopNameAr,
						StudentName:    order.CustomerName,
						StudentPhone:   order.CustomerPhone,
						TotalPrice:     item.Price,
						Status:         "pending",
						StatusAr:       "قيد الانتظار",
						EnrollmentDate: order.CreatedAt,
						CreatedAt:      time.Now(),
						UpdatedAt:      time.Now(),
					}
					
					err = enrollmentService.CreateEnrollment(enrollment)
					if err != nil {
						fmt.Printf("Warning: Failed to create enrollment: %v\n", err)
					}
				}
			}
		}
		fmt.Printf("Created order %s with status 'pending' - session attendance NOT incremented\n", order.ID.String())
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

	// Insert order items with session information
	itemQuery := `
		INSERT INTO order_items (
			id, order_id, workshop_id, session_id, run_id,
			workshop_name, workshop_name_ar, price, quantity, subtotal
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	for _, item := range order.Items {
		_, err = tx.Exec(itemQuery,
			item.ID, item.OrderID, item.WorkshopID, item.SessionID, item.RunID,
			item.WorkshopName, item.WorkshopNameAr, item.Price, item.Quantity, item.Subtotal,
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
		SELECT 
			oi.id, oi.order_id, oi.workshop_id, oi.session_id, oi.run_id,
			oi.workshop_name, oi.workshop_name_ar, oi.price, oi.quantity, oi.subtotal
		FROM order_items oi
		WHERE oi.order_id = $1
	`

	rows, err := database.Instance.Query(query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to query order items: %w", err)
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		var subtotal sql.NullFloat64
		
		err := rows.Scan(
			&item.ID, &item.OrderID, &item.WorkshopID, &item.SessionID, &item.RunID,
			&item.WorkshopName, &item.WorkshopNameAr, &item.Price, &item.Quantity, &subtotal,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}
		
		if subtotal.Valid {
			item.Subtotal = subtotal.Float64
		} else {
			item.Subtotal = item.Price * float64(item.Quantity)
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
	// First, get the current order to check what sessions need updating
	var currentOrder *models.Order
	currentStatus := ""
	
	// Try to get order from database first
	dbOrder, err := s.getOrderByID(orderID)
	if err == nil && dbOrder != nil {
		currentOrder = dbOrder
		currentStatus = dbOrder.Status
	} else {
		// Fall back to in-memory
		for _, order := range s.orders {
			if order.ID == orderID {
				currentOrder = &order
				currentStatus = order.Status
				break
			}
		}
	}
	
	if currentOrder == nil {
		return fmt.Errorf("order not found")
	}

	// Try to update in database first
	err = s.updateOrderStatusInDatabase(orderID, newStatus)
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

				break
			}
		}
	}

	// Handle session attendance updates only if status actually changed
	fmt.Printf("Order status change: %s -> %s for order %s\n", currentStatus, newStatus, orderID.String())
	
	// Only update attendance if the status is actually changing
	if currentStatus != newStatus {
		if currentStatus == "pending" && newStatus == "paid" {
			// Increment session attendance when order is paid
			fmt.Printf("Incrementing session attendance for %d items\n", len(currentOrder.Items))
			for _, item := range currentOrder.Items {
				if item.SessionID != nil {
					fmt.Printf("Updating session %s with quantity %d\n", item.SessionID.String(), item.Quantity)
					err := s.updateSessionAttendance(*item.SessionID, item.Quantity)
					if err != nil {
						fmt.Printf("Warning: Failed to increment session attendance for session %s: %v\n", item.SessionID.String(), err)
						// Don't return error - allow the status update to succeed even if attendance update fails
						// This prevents HTMX errors and ensures order status is properly updated
					}
				} else {
					fmt.Printf("Item has no session ID: %s\n", item.WorkshopName)
				}
			}
		} else if currentStatus == "paid" && (newStatus == "cancelled" || newStatus == "pending") {
			// Decrement session attendance when paid order is cancelled or reverted
			fmt.Printf("Decrementing session attendance for %d items\n", len(currentOrder.Items))
			for _, item := range currentOrder.Items {
				if item.SessionID != nil {
					fmt.Printf("Updating session %s with quantity -%d\n", item.SessionID.String(), item.Quantity)
					err := s.updateSessionAttendance(*item.SessionID, -item.Quantity)
					if err != nil {
						fmt.Printf("Warning: Failed to decrement session attendance for session %s: %v\n", item.SessionID.String(), err)
						// Don't return error - allow the status update to succeed even if attendance update fails
					}
				}
			}
		}
	} else {
		fmt.Printf("Status unchanged (%s), skipping attendance update\n", currentStatus)
	}
	
	return nil
}

func (s *OrderService) getOrderByID(orderID uuid.UUID) (*models.Order, error) {
	query := `
		SELECT o.id, o.creator_id, o.customer_name, o.customer_phone, 
		       o.total_amount, o.status, o.order_source, o.created_at, o.updated_at
		FROM orders o
		WHERE o.id = $1
	`
	
	var order models.Order
	err := database.Instance.QueryRow(query, orderID).Scan(
		&order.ID, &order.CreatorID, &order.CustomerName, &order.CustomerPhone,
		&order.TotalAmount, &order.Status, &order.OrderSource, &order.CreatedAt, &order.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get order by ID: %w", err)
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
	itemsQuery := `
		SELECT id, order_id, workshop_id, session_id, run_id,
		       workshop_name, workshop_name_ar, price, quantity, subtotal
		FROM order_items
		WHERE order_id = $1
	`
	
	rows, err := database.Instance.Query(itemsQuery, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}
	defer rows.Close()
	
	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		err := rows.Scan(
			&item.ID, &item.OrderID, &item.WorkshopID, &item.SessionID, &item.RunID,
			&item.WorkshopName, &item.WorkshopNameAr, &item.Price, &item.Quantity, &item.Subtotal,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}
		items = append(items, item)
	}
	
	order.Items = items
	return &order, nil
}

func (s *OrderService) updateSessionAttendance(sessionID uuid.UUID, quantityChange int) error {
	query := `
		UPDATE workshop_sessions
		SET current_attendees = current_attendees + $1
		WHERE id = $2
		AND (max_attendees = 0 OR current_attendees + $1 >= 0)
		AND (max_attendees = 0 OR current_attendees + $1 <= max_attendees OR $1 < 0)
	`
	
	result, err := database.Instance.Exec(query, quantityChange, sessionID)
	if err != nil {
		return fmt.Errorf("failed to update session attendance: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("session not found or attendance update would violate constraints")
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
	// First, get the order to check its status and items
	order, err := s.getOrderByID(orderID)
	if err != nil {
		return fmt.Errorf("failed to get order for deletion: %w", err)
	}
	if order == nil {
		return fmt.Errorf("order not found")
	}

	// If the order was paid, we need to decrement session attendance
	if order.Status == "paid" {
		fmt.Printf("Order %s is paid, decrementing session attendance for %d items\n", orderID.String(), len(order.Items))
		for _, item := range order.Items {
			if item.SessionID != nil {
				fmt.Printf("Decrementing session %s by %d seats\n", item.SessionID.String(), item.Quantity)
				err := s.updateSessionAttendance(*item.SessionID, -item.Quantity)
				if err != nil {
					fmt.Printf("Warning: Failed to decrement session attendance for session %s: %v\n", item.SessionID.String(), err)
					// Don't fail the deletion, just log the warning
				}
			}
		}
	}

	// Now delete from database
	err = s.deleteOrderFromDatabase(orderID)
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

// getSessionByID validates that a session exists and has a valid run_id
func (s *OrderService) getSessionByID(sessionID uuid.UUID) (*models.WorkshopSession, error) {
	// First, let's debug what we can find about this session
	debugQuery := `
		SELECT 
			ws.id, ws.workshop_id, ws.session_date, ws.start_time, ws.end_time,
			ws.max_attendees, ws.current_attendees, ws.run_id, ws.status,
			CASE WHEN ws.run_id IS NULL THEN 'NULL_RUN_ID'
				 WHEN NOT EXISTS (SELECT 1 FROM workshop_runs wr WHERE wr.id = ws.run_id) THEN 'INVALID_RUN_ID'
				 ELSE 'VALID_RUN_ID' END as run_id_status
		FROM workshop_sessions ws
		WHERE ws.id = $1
	`
	
	var session models.WorkshopSession
	var endTime sql.NullString
	var runIDStatus string
	
	err := database.Instance.QueryRow(debugQuery, sessionID).Scan(
		&session.ID, &session.WorkshopID, &session.SessionDate, 
		&session.StartTime, &endTime,
		&session.MaxAttendees, &session.CurrentAttendees, 
		&session.RunID, &session.Status, &runIDStatus,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session not found: %s", sessionID.String())
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query session: %w", err)
	}

	// Log debug information
	fmt.Printf("DEBUG: Session %s - Status: %s, RunIDStatus: %s, MaxAttendees: %d, CurrentAttendees: %d\n", 
		sessionID.String(), session.Status, runIDStatus, session.MaxAttendees, session.CurrentAttendees)

	// Check various conditions
	if session.Status == "cancelled" || session.Status == "completed" {
		return nil, fmt.Errorf("session is %s", session.Status)
	}
	
	if session.MaxAttendees > 0 && session.CurrentAttendees >= session.MaxAttendees {
		return nil, fmt.Errorf("session is full (%d/%d)", session.CurrentAttendees, session.MaxAttendees)
	}
	
	if runIDStatus == "NULL_RUN_ID" {
		return nil, fmt.Errorf("session has no run_id - needs data migration")
	}
	
	if runIDStatus == "INVALID_RUN_ID" {
		return nil, fmt.Errorf("session has invalid run_id - orphaned reference")
	}
	
	if endTime.Valid {
		session.EndTime = &endTime.String
	}
	
	return &session, nil
}
