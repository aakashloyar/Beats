package http

import (
	"net/http"
	"strings"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/tracks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			{
				h.ListTracks(w, r)
			}
		default:
			{
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		}

	})
	mux.HandleFunc("/tracks/", func(w http.ResponseWriter, r *http.Request) {

		path := r.URL.Path[len("/tracks/"):]
		parts := strings.Split(path, "/")

		trackID := parts[0]
		if trackID == "" {
			http.Error(w, "missing track id", http.StatusBadRequest)
			return
		}

		// /tracks/{id}
		if len(parts) == 1 {
			switch r.Method {
			case http.MethodGet:
				h.GetTrackByID(w, r, trackID)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}
		
		http.NotFound(w, r)
	})
}
