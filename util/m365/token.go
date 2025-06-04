package m365

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func RefreshToken(tenantID, refreshToken string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", os.Getenv("M365_CLIENT_ID"))
	data.Set("client_secret", os.Getenv("M365_CLIENT_SECRET"))
	data.Set("grant_type", "refresh_token")
	data.Set("scope", "https://graph.microsoft.com/.default")
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", "https://login.microsoftonline.com/"+tenantID+"/oauth2/v2.0/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}
	tokenResp.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	return &tokenResp, nil
}
