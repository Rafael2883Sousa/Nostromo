package controllers

import (
	"net/http"
	"github.com/gophish/gophish/util/m365"
	"github.com/gophish/gophish/models"
)

func M365AuthRedirect(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	redirectURI := config.Config.M365.RedirectURI
	clientID := config.Config.M365.ClientID

	authURL := "https://login.microsoftonline.com/" + tenantID + "/oauth2/v2.0/authorize" +
		"?client_id=" + clientID +
		"&response_type=code" +
		"&redirect_uri=" + url.QueryEscape(redirectURI) +
		"&response_mode=query" +
		"&scope=https%3A%2F%2Fgraph.microsoft.com%2F.default" +
		"&state=12345"

	http.Redirect(w, r, authURL, http.StatusFound)
}

func M365AuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	tenantID := r.URL.Query().Get("tenant_id")

	tokenResp, err := m365.ExchangeCodeForToken(r.Context(), tenantID, config.Config.M365.ClientID, config.Config.M365.ClientSecret, code, config.Config.M365.RedirectURI)
	if err != nil {
		http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
		return
	}

	tenant := models.M365Tenant{
		TenantID:     tenantID,
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
		ConsentedAt:  time.Now(),
	}

	err = models.DB.Save(&tenant).Error
	if err != nil {
		http.Error(w, "Failed to save tenant", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

func ImportGroupsFromGraph(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")

	var t models.M365Tenant
	if err := models.DB.Where("tenant_id = ?", tenantID).First(&t).Error; err != nil {
		http.Error(w, "Tenant not found", 404)
		return
	}

	groups, err := m365.FetchGroups(t.AccessToken)
	if err != nil {
		http.Error(w, "Failed to fetch groups", 500)
		return
	}

	for _, g := range groups {
		group := models.Group{
			Name:          g["displayName"].(string),
			ModifiedDate:  time.Now(),
			M365TenantID:  &t.ID,
		}
		models.DB.Create(&group)
	}
	http.Redirect(w, r, "/groups", 302)
}
