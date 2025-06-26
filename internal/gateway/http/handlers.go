package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/imirjar/rb-diver/internal/models"
)

type Service interface {
	ReportExecute(context.Context, string) (models.Data, error)
	ReportsList(context.Context, []string) ([]models.Report, error)
	// RoleList(context.Context, string) ([]models.Role, error)
}

// Use divers reports query fild as sql request in db
func (gw *HTTP) ReportExecute(w http.ResponseWriter, r *http.Request) {
	reportID := chi.URLParam(r, "id")

	result, err := gw.Service.ReportExecute(context.Background(), reportID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(&result); err != nil {
		log.Println("HANDLER ExecuteHandler Encode ERROR", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// ReportList consist of all diver reports
func (gw *HTTP) ReportsList(w http.ResponseWriter, r *http.Request) {
	roles := r.URL.Query()["role"]
	// log.Println(roles)

	result, err := gw.Service.ReportsList(r.Context(), roles)
	if err != nil {
		log.Println("HANDLER service ReportList ERROR", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(&result); err != nil {
		log.Println("HANDLER ExecuteHandler Encode ERROR", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// // ReportList consist of all diver reports
// func (gw *HTTP) RoleList(w http.ResponseWriter, r *http.Request) {
// 	reportID := r.URL.Query().Get("report_id")

// 	result, err := gw.Service.RoleList(r.Context(), reportID)
// 	if err != nil {
// 		log.Println("HANDLER service ReportList ERROR", err)
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	w.Header().Set("content-type", "application/json; charset=utf-8")
// 	w.WriteHeader(http.StatusOK)
// 	if err = json.NewEncoder(w).Encode(&result); err != nil {
// 		log.Println("HANDLER ExecuteHandler Encode ERROR", err)
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }

func (gw *HTTP) Info(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello! I,m a diver."))
}
