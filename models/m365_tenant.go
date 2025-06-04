package models

import "time"

type M365Tenant struct {
	ID           uint      `gorm:"primary_key"`
	TenantID     string    `gorm:"size:64;unique_index"`
	DisplayName  string
	ConsentedAt  time.Time
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}
