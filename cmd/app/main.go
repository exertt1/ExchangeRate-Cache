package main

import (
	"Excnahge-Cacher/api"
	"Excnahge-Cacher/config"
	core "Excnahge-Cacher/core/cache"
	"Excnahge-Cacher/handler"
)

type CurrencyLayerResponse struct {
	Success   bool   `json:"success"`
	Terms     string `json:"terms"`
	Privacy   string `json:"privacy"`
	Timestamp int64  `json:"timestamp"`
	Source    string `json:"source"`
}

func main() {
	cfg := config.NewConfig()

	apiHandler := api.NewAPIHandler(cfg)

	cache, err := core.NewCache(apiHandler)
	if err != nil {
		panic(err)
	}
	httpHandler := handler.NewHTTPHandler(cache)

}
