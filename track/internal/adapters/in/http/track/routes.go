package http

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/tracks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				h.CreateTrack(w, r)
			}
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
		switch r.Method {
		case http.MethodGet:
			{
				trackID := r.URL.Path[len("/tracks/"):]
				if trackID == "" {
					http.Error(w, "missing track id", http.StatusBadRequest)
					return
				}
				h.GetTrackByID(w, r, trackID)
			}
		default:
			{
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		}
	})
}
