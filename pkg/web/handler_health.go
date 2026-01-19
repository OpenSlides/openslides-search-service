package web

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

func healthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"healthy": true, "service":"search"}`)); err != nil {
			log.Errorf("error: writing response failed: %v\n", err)
		}
	}
}
