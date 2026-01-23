package track

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/tracks",func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost: {
			h.CreateTrack(w,r)
		}
		default: {
			http.Error(w,"method not allowed",http.StatusMethodNotAllowed)
		}
		}

	})
	mux.HandleFunc("/tracks/",func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet: {
			id := r.URL.Path[len("/tracks/"):]
			if id == "" {
				http.Error(w,"missing track id",http.StatusBadRequest)
				return 
			} 
			h.GetTrackById(w,r,id)
		}
		default: {
			http.Error(w,"method not allowed",http.StatusMethodNotAllowed)
		}
		}
	})
}