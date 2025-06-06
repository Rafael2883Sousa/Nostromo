package controllers

import (
	"net/http"
	"net/url"
	"time"
	"encoding/json"

	"github.com/gophish/gophish/config"
	"github.com/gophish/gophish/models"
	"github.com/gophish/gophish/util/m365"
)

func ListTenants(w http.ResponseWriter, r *http.Request) {
	tenants, err := models.GetAllTenants()
	if err != nil {
		http.Error(w, "Failed to fetch tenants", http.StatusInternalServerError)
		return
	}

	type tenantOut struct {
		ID          uint   `json:"id"`
		TenantID    string `json:"tenant_id"`
		DisplayName string `json:"display_name"`
	}

	var result []tenantOut
	for _, t := range tenants {
		result = append(result, tenantOut{
			ID:          t.ID,
			TenantID:    t.TenantID,
			DisplayName: t.DisplayName,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}


func M365AuthRedirect(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	redirectURI := config.Global.M365.RedirectURI
	clientID := config.Global.M365.ClientID

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

	log.Println("🔁 [OAuth] Callback recebido")
	log.Printf("🔁 [OAuth] Query params: code=%s tenant_id=%s\n", code, tenantID)

	tokenResp, err := m365.ExchangeCodeForToken(r.Context(), tenantID, config.Global.M365.ClientID, config.Global.M365.ClientSecret, code, config.Global.M365.RedirectURI)
	log.Println("➡️ Trocando código por token...")
	if err != nil {
		http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
		return
	}
	log.Println("➡️ Trocando código por token...")

	tenant := models.M365Tenant{
		TenantID:     tenantID,
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
		ConsentedAt:  time.Now(),
	}

	err = models.SaveTenant(&tenant)
	log.Println("📝 Salvando tenant no banco")
	if err != nil {
		http.Error(w, "Failed to save tenant", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusFound)
	log.Println("✅ Redirecionando para /groups")
}

func ImportGroupsFromGraph(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")

	t, err := models.GetTenantByID(tenantID)
	if err != nil {
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
		models.PutGroup(&group)
	}
	http.Redirect(w, r, "/groups", http.StatusFound)
}
