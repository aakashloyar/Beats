package http

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/artists", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				h.CreateArtist(w, r)
			}
		default:
			{
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		}
	})
	mux.HandleFunc("/artists/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			{
				artistID := r.URL.Path[len("/artists/"):]
				if artistID == "" {
					http.Error(w, "missing artist id", http.StatusBadRequest)
					return
				}
				h.GetArtistByID(w, r, artistID)
			}
		default:
			{
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		}
	})
}
