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

func SaveTenant(t *M365Tenant) error {
	return db.Save(t).Error
}

func GetTenantByID(tenantID string) (*M365Tenant, error) {
	var tenant M365Tenant
	err := db.Where("tenant_id = ?", tenantID).First(&tenant).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}


func GetAllTenants() ([]M365Tenant, error) {
	var tenants []M365Tenant
	err := db.Find(&tenants).Error
	return tenants, err
}