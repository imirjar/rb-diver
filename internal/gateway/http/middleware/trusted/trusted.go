package trusted

import (
	"log"
	"net/http"
)

func Middleware(addr string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Host != "addr" {
				log.Printf("Ratunku! I don't trust this host: %s!", r.Host)
				w.WriteHeader(http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
