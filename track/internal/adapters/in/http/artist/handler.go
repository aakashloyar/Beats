package http

import (
	"encoding/json"
	"github.com/aakashloyar/beats/track/internal/application/ports/in/artist"
	"net/http"
	"time"
)

type CreateArtistRequest struct {
	Name            string  `json:"name"`
	Bio             *string `json:"bio"`
	ProfileImageURL *string `json:"profile_image_url"`
}
type CreateTrackResponse struct {
	ArtistID string `json:"artistID"`
}

type GetArtistRequest struct {
	ArtistID string `json:"artistID"`
}

type GetArtistResponse struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Bio             *string   `json:"bio"`
	ProfileImageURL *string   `json:"profile_image_url"`
	CreatedAt       time.Time `json:"created_at"`
}

type Handler struct {
	createArtistService in.CreateArtistService
	getArtistService    in.GetArtistService
}

func NewHandler(createArtistService in.CreateArtistService, getArtistService in.GetArtistService) *Handler {
	return &Handler{
		createArtistService: createArtistService,
		getArtistService:    getArtistService,
	}
}

func (h *Handler) CreateArtist(w http.ResponseWriter, r *http.Request) {
	var req CreateArtistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}
	input := in.CreateArtistInput{
		Name:            req.Name,
		Bio:             req.Bio,
		ProfileImageURL: req.ProfileImageURL,
	}
	out, err := h.createArtistService.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp := CreateTrackResponse{
		ArtistID: out.ArtistID,
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) GetArtistByID(w http.ResponseWriter, r *http.Request, artistID string) {

	out, err := h.getArtistService.Execute(r.Context(), in.GetArtistInput{ArtistID: artistID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := GetArtistResponse{
		ID:              out.ID,
		Name:            out.Name,
		Bio:             out.Bio,
		ProfileImageURL: out.ProfileImageURL,
		CreatedAt:       out.CreatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}


