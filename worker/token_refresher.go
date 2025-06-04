package worker

import (
	"time"
	"github.com/gophish/gophish/models"
	"github.com/gophish/gophish/util/m365"
	log "github.com/gophish/gophish/logger"
)

func RefreshLoop() {
	for range time.Tick(10 * time.Minute) {
		tenants, err := models.GetAllTenants()
		if err != nil {
			log.Errorf("Failed to fetch tenants: %s", err)
			continue
		}

		for _, t := range tenants {
			if time.Until(t.ExpiresAt) < 5*time.Minute {
				newToken, err := m365.RefreshToken(t.TenantID, t.RefreshToken)
				if err != nil {
					log.Errorf("Error refreshing token for tenant %s: %s", t.TenantID, err)
					continue
				}
				t.AccessToken = newToken.AccessToken
				t.RefreshToken = newToken.RefreshToken
				t.ExpiresAt = newToken.ExpiresAt
				err = models.SaveTenant(&t)
				if err != nil {
					log.Errorf("Failed to save updated token for tenant %s: %s", t.TenantID, err)
				}

				log.Infof("Refreshed token for tenant %s", t.TenantID)
			}
		}
	}
}
