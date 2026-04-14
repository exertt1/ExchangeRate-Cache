package api

import (
	"Excnahge-Cacher/config"
	"net/http"
	"strings"
	"time"
)

const BaseURL = "https://api.exchangerate.host/live?access_key=YOUR_KEY"

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

func (h *APIHandler) GetAllCourses(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(h.PersonalURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	resp.
}
