package app

import (
	"context"
	"log"

	"github.com/imirjar/rb-diver/internal/gateway/http"
	"github.com/imirjar/rb-diver/internal/service"
	"github.com/imirjar/rb-diver/internal/storage"

	"github.com/imirjar/rb-diver/config"
)

// Run app
func Run(ctx context.Context) error {
	cfg := config.New()

	// Create storage layer for data
	storage := storage.New(cfg.DB)

	// Create service layer for app logic
	service := service.New()
	service.Storage = storage

	// Create HTTP server to serv http requests
	srv := http.New()
	srv.Service = service

	// Waiting for srv.Start's ending
	done := make(chan bool)
	go func() {

		// Start HTTP server
		err := srv.Start(ctx, cfg.Port, cfg.Michman)
		if err != nil {
			log.Print(err)
		}

		done <- true
	}()

	<-done
	return nil
}
