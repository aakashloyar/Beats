package http

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/audio-variants", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				h.CreateAudioVariant(w, r)
			}
		default:
			{
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		}

	})
}
