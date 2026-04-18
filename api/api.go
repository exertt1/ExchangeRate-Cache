package api

import (
	"Excnahge-Cacher/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const BaseURL = "https://api.exchangerate.host/live?access_key=YOUR_KEY"

type CurrencyLayerResponse struct {
	Success   bool               `json:"success"`
	Terms     string             `json:"terms"`
	Privacy   string             `json:"privacy"`
	Timestamp int64              `json:"timestamp"`
	Source    string             `json:"source"`
	Quotes    map[string]float64 `json:"quotes"`
}

type APIHandler struct {
	Config      *config.Config
	PersonalURL string
}

func NewAPIHandler(cfg *config.Config) *APIHandler {
	api := &APIHandler{Config: cfg}
	PersonalURL := strings.Replace(BaseURL, "YOUR_KEY", cfg.APIKey, 1)
	api.PersonalURL = PersonalURL
	return api
}

func (h *APIHandler) GetAllCourses() (*CurrencyLayerResponse, error) {
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(h.PersonalURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	result, err := ParseCurrencyLayerResponse(body)
	if err != nil {
		return nil, err
	}
	return result, err
}

func ParseCurrencyLayerResponse(body []byte) (*CurrencyLayerResponse, error) {
	var resp CurrencyLayerResponse
	err := json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("Bad response")
	}
	return &resp, nil
}
