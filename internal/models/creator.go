package models

import "time"

type Creator struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	NameAr    string    `json:"name_ar"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	Plan      string    `json:"plan"`
	PlanAr    string    `json:"plan_ar"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
