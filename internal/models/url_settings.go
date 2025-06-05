package models

import (
	"time"

	"github.com/google/uuid"
)

type URLSettings struct {
	ID          uuid.UUID  `json:"id"`
	CreatorID   uuid.UUID  `json:"creator_id"`
	Username    string     `json:"username"`
	ChangesUsed int        `json:"changes_used"`
	MaxChanges  int        `json:"max_changes"`
	LastChanged *time.Time `json:"last_changed"` // Changed to pointer to handle NULL
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type URLUpdateRequest struct {
	Username string `form:"username" json:"username"`
}

type URLValidationResult struct {
	IsValid        bool   `json:"is_valid"`
	ErrorMessage   string `json:"error_message"`
	ErrorMessageAr string `json:"error_message_ar"`
}
