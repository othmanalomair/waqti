// internal/models/workshop.go - Updated with missing fields
package models

import "time"

type Workshop struct {
	ID             int       `json:"id"`
	CreatorID      int       `json:"creator_id"`
	Name           string    `json:"name"` // Added missing Name field
	Title          string    `json:"title"`
	TitleAr        string    `json:"title_ar"`
	Description    string    `json:"description"`
	DescriptionAr  string    `json:"description_ar"`
	Price          float64   `json:"price"`
	Currency       string    `json:"currency"` // Added missing Currency field
	Duration       int       `json:"duration"` // in minutes
	MaxStudents    int       `json:"max_students"`
	Category       string    `json:"category"`
	CategoryAr     string    `json:"category_ar"`
	Status         string    `json:"status"` // Added missing Status field (draft, published, etc.)
	IsActive       bool      `json:"is_active"`
	IsFree         bool      `json:"is_free"`         // Added missing IsFree field
	IsRecurring    bool      `json:"is_recurring"`    // Added missing IsRecurring field
	RecurrenceType string    `json:"recurrence_type"` // Added missing RecurrenceType field (weekly, monthly)
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type WorkshopSession struct {
	ID          int       `json:"id"`
	WorkshopID  int       `json:"workshop_id"`
	Date        time.Time `json:"date"`
	StartTime   string    `json:"start_time"`
	Duration    float64   `json:"duration"` // in hours
	IsCompleted bool      `json:"is_completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type WorkshopImage struct {
	ID         int    `json:"id"`
	WorkshopID int    `json:"workshop_id"`
	ImageURL   string `json:"image_url"`
	IsCover    bool   `json:"is_cover"`
	SortOrder  int    `json:"sort_order"`
}
