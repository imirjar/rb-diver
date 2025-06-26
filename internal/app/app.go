package app

import (
	"context"
	"log"

	"github.com/imirjar/rb-diver/internal/gateway/amqp"
	"github.com/imirjar/rb-diver/internal/gateway/http"
	"github.com/imirjar/rb-diver/internal/service"
	"github.com/imirjar/rb-diver/internal/storage/reports"
	"github.com/imirjar/rb-diver/internal/storage/target"

	"github.com/imirjar/rb-diver/config"
	"golang.org/x/sync/errgroup"
)

// Run app
func Run(ctx context.Context) error {
	cfg := config.New()

	// Create storage layer for data
	targetStorage := target.New()
	err := targetStorage.Connect(ctx, cfg.DB)
	if err != nil {
		return err
	}
	defer targetStorage.Disconnect()

	rs := reports.New()
	err = rs.Connect(ctx, cfg.Mongo)
	if err != nil {
		return err
	}
	defer rs.Disconnect()

	// Create service layer for app logic
	svc := service.New()
	svc.Storage = targetStorage
	svc.RS = rs

	// Create HTTP server to serv http requests
	srv := http.New()
	srv.Service = svc

	amqpServer := amqp.New()
	amqpServer.Connect(ctx, cfg.Rabbit)
	defer amqpServer.Disconnect()

	amqpServer.Service = svc

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		// Start HTTP server
		if err := srv.Start(ctx, cfg.Port); err != nil {
			log.Printf("HTTP server error: %v", err)
			return err
		}
		return nil
	})

	g.Go(func() error {
		// Start AMQP server
		if err := amqpServer.Start(ctx, cfg.Rabbit); err != nil {
			log.Printf("AMQP server error: %v", err)
			return err
		}
		return nil
	})

	// Wait for all servers to finish or one to error
	return g.Wait()
}
