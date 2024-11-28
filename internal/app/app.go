package app

import (
	"context"
	"fmt"
	"log"

	"github.com/imirjar/rb-diver/internal/gateway/http"
	"github.com/imirjar/rb-diver/internal/models"
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
		// registration in Michman
		diver := models.Diver{
			Name: cfg.Name,
			Addr: fmt.Sprintf("%s:%s", cfg.Addr, cfg.Port),
			Port: cfg.Port,

			Michman: cfg.Michman,
		}

		err := srv.Registrate(ctx, diver)
		if err != nil {
			log.Print(err)
			done <- true
		}

		// Start HTTP server
		err = srv.Start(ctx, cfg.Port, cfg.Michman)
		if err != nil {
			panic(err)
		}

		done <- true
	}()

	<-done
	return nil
}
