package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"waqti/internal/models"
	"waqti/internal/services"
	"waqti/web/templates"

	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	creatorService *services.CreatorService
	orderService   *services.OrderService
}

func NewOrderHandler(creatorService *services.CreatorService, orderService *services.OrderService) *OrderHandler {
	return &OrderHandler{
		creatorService: creatorService,
		orderService:   orderService,
	}
}

// CreateOrder handles order creation from the store page
func (h *OrderHandler) CreateOrder(c echo.Context) error {
	var request models.CreateOrderRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Validate required fields
	if request.CustomerName == "" || request.CustomerPhone == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Customer name and phone are required"})
	}

	if len(request.Items) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "At least one item is required"})
	}

	// For demo purposes, use creator ID 1
	// In real implementation, you would extract this from the URL or auth context
	creatorID := 1

	order, err := h.orderService.CreateOrder(creatorID, request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Return success response
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"order":   order,
		"message": "Order created successfully",
	})
}

// ShowOrderTracking displays the order tracking page
func (h *OrderHandler) ShowOrderTracking(c echo.Context) error {
	// Get language from context
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	// Get query parameters for filtering
	timeRange := c.QueryParam("time_range")
	if timeRange == "" {
		timeRange = "days"
	}

	orderBy := c.QueryParam("order_by")
	if orderBy == "" {
		orderBy = "date"
	}

	orderDir := c.QueryParam("order_dir")
	if orderDir == "" {
		orderDir = "desc"
	}

	filter := models.EnrollmentFilter{
		TimeRange: timeRange,
		OrderBy:   orderBy,
		OrderDir:  orderDir,
	}

	// Get orders and stats
	orders := h.orderService.GetOrdersByCreatorID(1, filter)
	stats := h.orderService.GetOrderStats(1, timeRange)

	// Render template
	component := templates.OrderTrackingPage(orders, stats, filter, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

// FilterOrders handles HTMX filtering requests
func (h *OrderHandler) FilterOrders(c echo.Context) error {
	// Get language from context
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	// Get form parameters
	timeRange := c.FormValue("time_range")
	if timeRange == "" {
		timeRange = "days"
	}

	orderBy := c.FormValue("order_by")
	if orderBy == "" {
		orderBy = "date"
	}

	orderDir := c.FormValue("order_dir")
	if orderDir == "" {
		orderDir = "desc"
	}

	filter := models.EnrollmentFilter{
		TimeRange: timeRange,
		OrderBy:   orderBy,
		OrderDir:  orderDir,
	}

	// Get filtered data
	orders := h.orderService.GetOrdersByCreatorID(1, filter)
	stats := h.orderService.GetOrderStats(1, timeRange)

	// Return updated content via HTMX
	component := templates.OrderContent(orders, stats, filter, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

// UpdateOrderStatus handles individual order status updates
func (h *OrderHandler) UpdateOrderStatus(c echo.Context) error {
	orderID, err := strconv.Atoi(c.FormValue("order_id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid order ID")
	}

	newStatus := c.FormValue("status")
	if newStatus == "" {
		return c.String(http.StatusBadRequest, "Status is required")
	}

	// Validate status
	validStatuses := map[string]bool{
		"pending":   true,
		"paid":      true,
		"cancelled": true,
	}

	if !validStatuses[newStatus] {
		return c.String(http.StatusBadRequest, "Invalid status")
	}

	err = h.orderService.UpdateOrderStatus(orderID, newStatus)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error updating order status")
	}

	// Return updated order list via HTMX
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	// Get current filter from form or use defaults
	timeRange := c.FormValue("current_time_range")
	if timeRange == "" {
		timeRange = "days"
	}
	orderBy := c.FormValue("current_order_by")
	if orderBy == "" {
		orderBy = "date"
	}
	orderDir := c.FormValue("current_order_dir")
	if orderDir == "" {
		orderDir = "desc"
	}

	filter := models.EnrollmentFilter{
		TimeRange: timeRange,
		OrderBy:   orderBy,
		OrderDir:  orderDir,
	}

	orders := h.orderService.GetOrdersByCreatorID(1, filter)
	stats := h.orderService.GetOrderStats(1, timeRange)

	component := templates.OrderContent(orders, stats, filter, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

// DeleteOrder handles order deletion
func (h *OrderHandler) DeleteOrder(c echo.Context) error {
	orderID, err := strconv.Atoi(c.FormValue("order_id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid order ID")
	}

	err = h.orderService.DeleteOrder(orderID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error deleting order")
	}

	// Return updated list via HTMX
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	// Get current filter from form or use defaults
	timeRange := c.FormValue("current_time_range")
	if timeRange == "" {
		timeRange = "days"
	}
	orderBy := c.FormValue("current_order_by")
	if orderBy == "" {
		orderBy = "date"
	}
	orderDir := c.FormValue("current_order_dir")
	if orderDir == "" {
		orderDir = "desc"
	}

	filter := models.EnrollmentFilter{
		TimeRange: timeRange,
		OrderBy:   orderBy,
		OrderDir:  orderDir,
	}

	orders := h.orderService.GetOrdersByCreatorID(1, filter)
	stats := h.orderService.GetOrderStats(1, timeRange)

	component := templates.OrderContent(orders, stats, filter, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

// BulkAction handles bulk operations on orders
func (h *OrderHandler) BulkAction(c echo.Context) error {
	action := c.FormValue("action")
	status := c.FormValue("status")

	if action == "" {
		return c.String(http.StatusBadRequest, "Action is required")
	}

	// Validate action
	validActions := map[string]bool{
		"mark_paid": true,
		"cancel":    true,
	}

	if !validActions[action] {
		return c.String(http.StatusBadRequest, "Invalid action")
	}

	// Get all orders with the specified status
	orders := h.orderService.GetOrdersByCreatorID(1, models.EnrollmentFilter{
		TimeRange: "days",
		OrderBy:   "date",
		OrderDir:  "desc",
	})

	var updatedCount int
	for _, order := range orders {
		// Apply bulk action based on current status and action
		if status == "" || order.Status == status {
			switch action {
			case "mark_paid":
				if order.Status == "pending" {
					err := h.orderService.UpdateOrderStatus(order.ID, "paid")
					if err == nil {
						updatedCount++
					}
				}
			case "cancel":
				if order.Status == "pending" {
					err := h.orderService.UpdateOrderStatus(order.ID, "cancelled")
					if err == nil {
						updatedCount++
					}
				}
			}
		}
	}

	// Return updated list via HTMX
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	filter := models.EnrollmentFilter{
		TimeRange: "days",
		OrderBy:   "date",
		OrderDir:  "desc",
	}

	orders = h.orderService.GetOrdersByCreatorID(1, filter)
	stats := h.orderService.GetOrderStats(1, "days")

	component := templates.OrderContent(orders, stats, filter, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

// GetOrderDetails returns detailed information about a specific order (API endpoint)
func (h *OrderHandler) GetOrderDetails(c echo.Context) error {
	orderIDStr := c.Param("id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid order ID"})
	}

	// Get order details
	orders := h.orderService.GetOrdersByCreatorID(1, models.EnrollmentFilter{
		TimeRange: "all",
		OrderBy:   "date",
		OrderDir:  "desc",
	})

	var foundOrder *models.Order
	for _, order := range orders {
		if order.ID == orderID {
			foundOrder = &order
			break
		}
	}

	if foundOrder == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Order not found"})
	}

	return c.JSON(http.StatusOK, foundOrder)
}

// GetOrderStats returns order statistics (API endpoint)
func (h *OrderHandler) GetOrderStats(c echo.Context) error {
	timeRange := c.QueryParam("time_range")
	if timeRange == "" {
		timeRange = "days"
	}

	stats := h.orderService.GetOrderStats(1, timeRange)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"stats":      stats,
		"time_range": timeRange,
	})
}

// ExportOrders exports orders to CSV (future enhancement)
func (h *OrderHandler) ExportOrders(c echo.Context) error {
	// Get all orders
	orders := h.orderService.GetOrdersByCreatorID(1, models.EnrollmentFilter{
		TimeRange: "all",
		OrderBy:   "date",
		OrderDir:  "desc",
	})

	// Set headers for CSV download
	c.Response().Header().Set("Content-Type", "text/csv")
	c.Response().Header().Set("Content-Disposition", "attachment; filename=orders.csv")

	// Write CSV header
	csvContent := "Order ID,Customer Name,Customer Phone,Total Amount,Status,Order Date,Items\n"

	// Write order data
	for _, order := range orders {
		itemsStr := ""
		for i, item := range order.Items {
			if i > 0 {
				itemsStr += "; "
			}
			itemsStr += fmt.Sprintf("%s (%.2f KD)", item.WorkshopName, item.Price)
		}

		csvContent += fmt.Sprintf("%d,%s,%s,%.2f,%s,%s,\"%s\"\n",
			order.ID,
			order.CustomerName,
			order.CustomerPhone,
			order.TotalAmount,
			order.Status,
			order.CreatedAt.Format("2006-01-02 15:04:05"),
			itemsStr,
		)
	}

	return c.String(http.StatusOK, csvContent)
}

// MarkOrderAsViewed marks an order as viewed by the creator (for notification management)
func (h *OrderHandler) MarkOrderAsViewed(c echo.Context) error {
	orderID, err := strconv.Atoi(c.FormValue("order_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid order ID"})
	}

	// In a real implementation, you would update a "viewed" field in the database
	// For now, we'll just return success
	_ = orderID

	return c.JSON(http.StatusOK, map[string]string{"message": "Order marked as viewed"})
}
