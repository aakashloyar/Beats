package http

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/albums", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				h.CreateAlbum(w, r)
			}
		case http.MethodGet:
		{
			h.ListAlbums(w, r)
		}	
		default:
			{
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		}
	})
	mux.HandleFunc("/albums/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			{
				albumID := r.URL.Path[len("/albums/"):]
				if albumID == "" {
					http.Error(w, "missing album id", http.StatusBadRequest)
					return
				}
				h.GetAlbumByID(w, r, albumID)
			}
		default:
			{
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		}
	})
}
