package models

import (
	"time"

	"github.com/google/uuid"
)

type Workshop struct {
	ID             uuid.UUID `json:"id"`
	CreatorID      uuid.UUID `json:"creator_id"`
	Name           string    `json:"name"`
	Title          string    `json:"title"`
	TitleAr        string    `json:"title_ar"`
	Description    string    `json:"description"`
	DescriptionAr  string    `json:"description_ar"`
	Price          float64   `json:"price"`
	Currency       string    `json:"currency"`
	Duration       int       `json:"duration"` // in minutes
	MaxStudents    int       `json:"max_students"`
	Category       string    `json:"category"`
	CategoryAr     string    `json:"category_ar"`
	Status         string    `json:"status"`
	IsActive       bool      `json:"is_active"`
	IsFree         bool      `json:"is_free"`
	IsRecurring    bool      `json:"is_recurring"`
	RecurrenceType string    `json:"recurrence_type"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type WorkshopSession struct {
	ID          uuid.UUID `json:"id"`
	WorkshopID  uuid.UUID `json:"workshop_id"`
	Date        time.Time `json:"date"`
	StartTime   string    `json:"start_time"`
	Duration    float64   `json:"duration"` // in hours
	IsCompleted bool      `json:"is_completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type WorkshopImage struct {
	ID         uuid.UUID `json:"id"`
	WorkshopID uuid.UUID `json:"workshop_id"`
	ImageURL   string    `json:"image_url"`
	IsCover    bool      `json:"is_cover"`
	SortOrder  int       `json:"sort_order"`
}
