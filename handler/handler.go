package handler

import (
	core "Excnahge-Cacher/core/cache"
	"encoding/json"
	"net/http"
	"strings"
)

type HTTPHandler struct {
	Cache *core.Cache
}

func NewHTTPHandler(cache *core.Cache) *HTTPHandler {
	return &HTTPHandler{Cache: cache}
}

type AllCurrencyRequest struct {
}
type CurrenciesResponse struct {
	Currencies map[string]float64 `json:"currencies"`
}

func (h *HTTPHandler) GetAllCurrency(w http.ResponseWriter, r *http.Request) {
	cacheItems := h.Cache.GetAll()
	currencies := make(map[string]float64)
	for _, val := range cacheItems {
		currencies[val.FromValue+val.ToValue] = val.Currency
	}
	var resp CurrenciesResponse
	resp.Currencies = currencies
	err := json.NewEncoder(w).Encode(&resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *HTTPHandler) GetCurrency(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	if from == "" || to == "" {
		h.GetAllCurrency(w, r)
	}
	from = strings.ToUpper(from)
	to = strings.ToUpper(to)
	cacheItem := h.Cache.Get(from + to)
	currency := make(map[string]float64)
	currency[from+to] = cacheItem.Currency
	var resp CurrenciesResponse
	resp.Currencies = currency
	err := json.NewEncoder(w).Encode(&resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
