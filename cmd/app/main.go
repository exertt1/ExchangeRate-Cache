package main

import (
	"Excnahge-Cacher/api"
	"Excnahge-Cacher/config"
	core "Excnahge-Cacher/core/cache"
	"Excnahge-Cacher/handler"
	"Excnahge-Cacher/router"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	cacheRouter := router.NewRouter(httpHandler)
	srv := &http.Server{
		Addr:    "localhost:8888",
		Handler: cacheRouter,
	}
	go func() {
		log.Println("Server is started")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")

}
