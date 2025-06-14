package models

import (
	"time"
	"encoding/hex"
	"crypto/rand"
)

type M365Tenant struct {
	ID           string    `gorm:"primary_key" json:"id"`
	Name         string    `json:"tenant_name"`
	TenantID     string    `json:"tenant_id"`
	ClientID     string    `json:"client_id"`
	ClientSecret string    `json:"client_secret"`
	CreatedAt    time.Time `json:"created_at"`
}

func GetAllTenants() ([]M365Tenant, error) {
	var tenants []M365Tenant
	err := db.Find(&tenants).Error
	return tenants, err
}
 
func SaveTenant(t *M365Tenant) error {
	if t.ID == "" {
		t.ID = GenerateSecureRandomString(12)
	}
	return db.Create(t).Error
}

func GetTenantByID(id string) (*M365Tenant, error) {
	var tenant M365Tenant
	err := db.Where("id = ?", id).First(&tenant).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func DeleteTenantByID(id string) error {
	return db.Where("id = ?", id).Delete(&M365Tenant{}).Error
}

func GenerateSecureRandomString(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err) // Falha catastrófica: não consegue gerar entropia
	}
	return hex.EncodeToString(bytes)[:length]
}


