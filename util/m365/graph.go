package m365

import (
	"encoding/json"
	"net/http"
)

func FetchGroups(accessToken string) ([]map[string]interface{}, error) {
	req, _ := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/groups?$select=id,displayName", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Value []map[string]interface{} `json:"value"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result.Value, err
}
