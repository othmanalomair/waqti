package handlers

import (
	"fmt"
	"net/http"
	"waqti/internal/database"
	"waqti/internal/middleware"
	"waqti/internal/models"
	"waqti/internal/services"
	"waqti/web/templates"

	"github.com/google/uuid"

	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	creatorService *services.CreatorService
	orderService   *services.OrderService
}

// Helper function to load shop settings
func (h *OrderHandler) loadShopSettings(creatorID uuid.UUID) *models.ShopSettings {
	dbShopSettings, err := database.Instance.GetShopSettingsByCreatorID(creatorID)
	if err != nil || dbShopSettings == nil {
		return nil
	}

	// Helper function to safely dereference string pointers
	getStringValue := func(s *string) string {
		if s != nil {
			return *s
		}
		return ""
	}
	
	return &models.ShopSettings{
		ID:                 dbShopSettings.ID,
		CreatorID:          dbShopSettings.CreatorID,
		LogoURL:            getStringValue(dbShopSettings.LogoURL),
		CreatorName:        getStringValue(dbShopSettings.CreatorName),
		CreatorNameAr:      getStringValue(dbShopSettings.CreatorNameAr),
		SubHeader:          getStringValue(dbShopSettings.SubHeader),
		SubHeaderAr:        getStringValue(dbShopSettings.SubHeaderAr),
		EnrollmentWhatsApp: getStringValue(dbShopSettings.EnrollmentWhatsApp),
		ContactWhatsApp:    getStringValue(dbShopSettings.ContactWhatsApp),
		CheckoutLanguage:   dbShopSettings.CheckoutLanguage,
		GreetingMessage:    getStringValue(dbShopSettings.GreetingMessage),
		GreetingMessageAr:  getStringValue(dbShopSettings.GreetingMessageAr),
		CurrencySymbol:     dbShopSettings.CurrencySymbol,
		CurrencySymbolAr:   dbShopSettings.CurrencySymbolAr,
		CreatedAt:          dbShopSettings.CreatedAt,
		UpdatedAt:          dbShopSettings.UpdatedAt,
	}
}

func NewOrderHandler(creatorService *services.CreatorService, orderService *services.OrderService) *OrderHandler {
	return &OrderHandler{
		creatorService: creatorService,
		orderService:   orderService,
	}
}

func (h *OrderHandler) CreateOrder(c echo.Context) error {
	var request models.CreateOrderRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Debug logging
	fmt.Printf("Order creation request: Creator=%s, Customer=%s, Phone=%s, Items=%d\n", 
		request.CreatorUsername, request.CustomerName, request.CustomerPhone, len(request.Items))

	if request.CustomerName == "" || request.CustomerPhone == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Customer name and phone are required"})
	}

	if len(request.Items) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "At least one item is required"})
	}

	// Look up creator by username if provided
	var creatorID uuid.UUID
	if request.CreatorUsername != "" {
		// Use database lookup directly (same as store page handler)
		creator, err := database.Instance.GetCreatorByUsername(request.CreatorUsername)
		if err != nil || creator == nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": fmt.Sprintf("Creator not found for username: %s", request.CreatorUsername),
			})
		}
		creatorID = creator.ID
	} else {
		// Fallback to default creator for demo
		creatorID = uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	}

	order, err := h.orderService.CreateOrder(creatorID, request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"order":   order,
		"message": "Order created successfully",
	})
}

func (h *OrderHandler) ShowOrderTracking(c echo.Context) error {
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	// Get current authenticated creator
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		// This should be handled by middleware, but just in case
		return c.Redirect(http.StatusSeeOther, "/signin")
	}

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

	creatorID := dbCreator.ID
	orders := h.orderService.GetOrdersByCreatorID(creatorID, filter)
	stats := h.orderService.GetOrderStats(creatorID, timeRange)

	// Get shop settings
	dbShopSettings, err := database.Instance.GetShopSettingsByCreatorID(dbCreator.ID)
	var shopSettings *models.ShopSettings
	if err == nil && dbShopSettings != nil {
		// Helper function to safely dereference string pointers
		getStringValue := func(s *string) string {
			if s != nil {
				return *s
			}
			return ""
		}
		
		shopSettings = &models.ShopSettings{
			ID:                 dbShopSettings.ID,
			CreatorID:          dbShopSettings.CreatorID,
			LogoURL:            getStringValue(dbShopSettings.LogoURL),
			CreatorName:        getStringValue(dbShopSettings.CreatorName),
			CreatorNameAr:      getStringValue(dbShopSettings.CreatorNameAr),
			SubHeader:          getStringValue(dbShopSettings.SubHeader),
			SubHeaderAr:        getStringValue(dbShopSettings.SubHeaderAr),
			EnrollmentWhatsApp: getStringValue(dbShopSettings.EnrollmentWhatsApp),
			ContactWhatsApp:    getStringValue(dbShopSettings.ContactWhatsApp),
			CheckoutLanguage:   dbShopSettings.CheckoutLanguage,
			GreetingMessage:    getStringValue(dbShopSettings.GreetingMessage),
			GreetingMessageAr:  getStringValue(dbShopSettings.GreetingMessageAr),
			CurrencySymbol:     dbShopSettings.CurrencySymbol,
			CurrencySymbolAr:   dbShopSettings.CurrencySymbolAr,
			CreatedAt:          dbShopSettings.CreatedAt,
			UpdatedAt:          dbShopSettings.UpdatedAt,
		}
	}

	component := templates.OrderTrackingPage(orders, stats, filter, shopSettings, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *OrderHandler) FilterOrders(c echo.Context) error {
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	// Get current authenticated creator
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}

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

	creatorID := dbCreator.ID
	orders := h.orderService.GetOrdersByCreatorID(creatorID, filter)
	stats := h.orderService.GetOrderStats(creatorID, timeRange)

	shopSettings := h.loadShopSettings(dbCreator.ID)
	component := templates.OrderContent(orders, stats, filter, shopSettings, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *OrderHandler) UpdateOrderStatus(c echo.Context) error {
	orderIDStr := c.FormValue("order_id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid order ID")
	}

	newStatus := c.FormValue("status")
	if newStatus == "" {
		return c.String(http.StatusBadRequest, "Status is required")
	}

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

	// Get current authenticated creator
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}

	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

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

	creatorID := dbCreator.ID
	orders := h.orderService.GetOrdersByCreatorID(creatorID, filter)
	stats := h.orderService.GetOrderStats(creatorID, timeRange)

	shopSettings := h.loadShopSettings(dbCreator.ID)
	component := templates.OrderContent(orders, stats, filter, shopSettings, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *OrderHandler) DeleteOrder(c echo.Context) error {
	orderIDStr := c.FormValue("order_id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid order ID")
	}

	err = h.orderService.DeleteOrder(orderID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error deleting order")
	}

	// Get current authenticated creator
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}

	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

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

	creatorID := dbCreator.ID
	orders := h.orderService.GetOrdersByCreatorID(creatorID, filter)
	stats := h.orderService.GetOrderStats(creatorID, timeRange)

	shopSettings := h.loadShopSettings(dbCreator.ID)
	component := templates.OrderContent(orders, stats, filter, shopSettings, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *OrderHandler) BulkAction(c echo.Context) error {
	action := c.FormValue("action")
	status := c.FormValue("status")

	if action == "" {
		return c.String(http.StatusBadRequest, "Action is required")
	}

	validActions := map[string]bool{
		"mark_paid": true,
		"cancel":    true,
	}

	if !validActions[action] {
		return c.String(http.StatusBadRequest, "Invalid action")
	}

	// Get current authenticated creator
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}

	creatorID := dbCreator.ID
	orders := h.orderService.GetOrdersByCreatorID(creatorID, models.EnrollmentFilter{
		TimeRange: "days",
		OrderBy:   "date",
		OrderDir:  "desc",
	})

	var updatedCount int
	for _, order := range orders {
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

	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	filter := models.EnrollmentFilter{
		TimeRange: "days",
		OrderBy:   "date",
		OrderDir:  "desc",
	}

	orders = h.orderService.GetOrdersByCreatorID(creatorID, filter)
	stats := h.orderService.GetOrderStats(creatorID, "days")

	shopSettings := h.loadShopSettings(dbCreator.ID)
	component := templates.OrderContent(orders, stats, filter, shopSettings, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *OrderHandler) GetOrderDetails(c echo.Context) error {
	orderIDStr := c.Param("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid order ID"})
	}

	creatorID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	orders := h.orderService.GetOrdersByCreatorID(creatorID, models.EnrollmentFilter{
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

func (h *OrderHandler) GetOrderStats(c echo.Context) error {
	timeRange := c.QueryParam("time_range")
	if timeRange == "" {
		timeRange = "days"
	}

	creatorID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	stats := h.orderService.GetOrderStats(creatorID, timeRange)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"stats":      stats,
		"time_range": timeRange,
	})
}

func (h *OrderHandler) ExportOrders(c echo.Context) error {
	creatorID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	orders := h.orderService.GetOrdersByCreatorID(creatorID, models.EnrollmentFilter{
		TimeRange: "all",
		OrderBy:   "date",
		OrderDir:  "desc",
	})

	c.Response().Header().Set("Content-Type", "text/csv")
	c.Response().Header().Set("Content-Disposition", "attachment; filename=orders.csv")

	csvContent := "Order ID,Customer Name,Customer Phone,Total Amount,Status,Order Date,Items\n"

	for _, order := range orders {
		itemsStr := ""
		for i, item := range order.Items {
			if i > 0 {
				itemsStr += "; "
			}
			itemsStr += fmt.Sprintf("%s (%.2f KD)", item.WorkshopName, item.Price)
		}

		csvContent += fmt.Sprintf("%s,%s,%s,%.2f,%s,%s,\"%s\"\n",
			order.ID.String(),
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

func (h *OrderHandler) MarkOrderAsViewed(c echo.Context) error {
	orderIDStr := c.FormValue("order_id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid order ID"})
	}

	_ = orderID

	return c.JSON(http.StatusOK, map[string]string{"message": "Order marked as viewed"})
}
