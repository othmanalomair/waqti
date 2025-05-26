package models

import "time"

type URLSettings struct {
	ID          int       `json:"id"`
	CreatorID   int       `json:"creator_id"`
	Username    string    `json:"username"`
	ChangesUsed int       `json:"changes_used"`
	MaxChanges  int       `json:"max_changes"`
	LastChanged time.Time `json:"last_changed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type URLUpdateRequest struct {
	Username string `form:"username" json:"username"`
}

type URLValidationResult struct {
	IsValid        bool   `json:"is_valid"`
	ErrorMessage   string `json:"error_message"`
	ErrorMessageAr string `json:"error_message_ar"`
}
