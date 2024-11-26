package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/imirjar/rb-diver/internal/models"
	"github.com/imirjar/rb-diver/internal/service"
)

type Service interface {
	Execute(ctx context.Context, id string) ([]map[string]interface{}, error)
	ReportsList(ctx context.Context) (string, error)
}

type HTTP struct {
	client  *client
	Service Service
}

func New() *HTTP {
	return &HTTP{
		client:  &client{http.DefaultClient},
		Service: service.New(),
	}
}

func (gw *HTTP) Start(ctx context.Context, addr string, michman string) error {
	router := chi.NewRouter()

	router.Get("/", gw.Info)

	router.Route("/reports", func(update chi.Router) {
		update.Get("/", gw.ReportsList)
		update.Get("/generate/{id}", gw.GenerateReport)
	})

	//for new usecases add new route
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", addr),
		Handler: router,
	}

	log.Printf("Run app on PORT %s", addr)
	return srv.ListenAndServe()
}

func (gw *HTTP) Registrate(ctx context.Context, addr string) error {
	log.Print("Michman!!! I'm here! Under the water!")

	diver := models.Diver{
		Name: "MyMac",
		IP:   "192.168.0.1",
	}

	md, err := json.Marshal(diver)
	if err != nil {
		log.Print(err)
	}

	reader := bytes.NewReader(md)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://"+addr+"/connect", reader)
	if err != nil {
		log.Print(err)
		return err
	}
	req.Header.Add("X-Real-IP", addr)
	status, err := gw.client.POST(req)
	if err != nil {
		log.Print(err)
		return err
	}
	log.Print(status)
	return nil
}
