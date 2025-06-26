package http

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

type HTTP struct {
	router  chi.Mux
	Service Service
}

func New() *HTTP {
	router := chi.NewRouter()
	return &HTTP{
		router: *router,
	}

}

func (gw *HTTP) Start(ctx context.Context, addr, michman string) error {

	gw.router.Get("/", gw.Info)

	// gw.router.Route("/roles", func(roles chi.Router) {
	// 	roles.Get("/", gw.RoleList)
	// })

	gw.router.Route("/reports", func(reports chi.Router) {
		// reports.Get("/has_role", gw.GetReport)
		reports.Post("/{id}", gw.ReportExecute)
		// reports.Put("/{id}", gw.ReportUpdate)
		reports.Get("/", gw.ReportsList)
	})

	//for new usecases add new route
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", addr),
		Handler: &gw.router,
	}

	log.Printf("Run app on PORT %s", addr)
	return srv.ListenAndServe()
}
