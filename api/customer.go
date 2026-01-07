package api

import (
	"encoding/json"
	"fmt"
	"github.com/AndersKaae/legaldesk_psp_sync/config"
	"io"
	"net/http"
)

func GetCustomer(customerId string) (map[string]any, error) {
	url := "https://api.frisbii.com/v1/invoice/" + customerId

	cfg := config.LoadConfig()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.SetBasicAuth(cfg.Psp_api_key_dk, "")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return map[string]any{}, fmt.Errorf("customer %s not found", customerId)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	// Pretty-print so you can inspect it easily
	pretty, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(pretty))

	return data, nil
}
