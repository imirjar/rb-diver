package http

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/imirjar/rb-diver/internal/service"
)

type Service interface {
	Execute(ctx context.Context, id string) (string, error)
	ReportsList(ctx context.Context) (string, error)
}

type HTTP struct {
	Service Service
}

type Config interface {
	GetDiverAddr() string
	GetMichmanAddr() string //allow req only for this addr
	GetSecret() string
}

func New() *HTTP {
	return &HTTP{
		Service: service.New(),
	}
}

func (a *HTTP) Start(ctx context.Context, addr string) error {
	router := chi.NewRouter()

	// router.Use(middleware.Encryptor(a.config.GetSecret()))

	router.Route("/reports", func(update chi.Router) {
		update.Post("/execute/", a.ExecuteHandler)
		update.Post("/list/", a.ReportsListHandler)
	})

	//for new usecases add new route
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	log.Printf("Run app on PORT %s", addr)
	return srv.ListenAndServe()
}
