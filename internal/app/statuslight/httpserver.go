package statuslight

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// HTTPServer is a HTTP server processing status light commands.
type HTTPServer struct {
	port        int
	statusLight *StatusLight
}

// NewHTTPServer returns initialized HTTPServer object.
func NewHTTPServer(port int, statusLight *StatusLight) *HTTPServer {
	return &HTTPServer{
		port:        port,
		statusLight: statusLight,
	}
}

// ListenAndServe starts HTTP server.
func (s *HTTPServer) ListenAndServe() error {
	r := mux.NewRouter()
	v1 := r.PathPrefix("/api/v1/").Subrouter()

	v1.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		statusHandler(w, r, s.statusLight)
	}).Methods("POST")

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%d", s.port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv.ListenAndServe()
}

// statusHandler processes HTTP API calls.
func statusHandler(w http.ResponseWriter, r *http.Request, statusLight *StatusLight) {
	var s status
	if r.Body == nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err = statusLight.processStatus(s)
	if err != nil {
		log.Printf("processStatus error: %s\n", err)
		http.Error(w, "statuslight error", http.StatusInternalServerError)
	}
}
