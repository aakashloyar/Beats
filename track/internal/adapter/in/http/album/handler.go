package http

import (
	"encoding/json"
	"github.com/aakashloyar/beats/track/internal/application/ports/in/album"
	"net/http"
	"time"
)

type CreateAlbumRequest struct {
	Title         string     `json:"title"`
	CoverImageURL *string    `json:"cover_image_url"`
	ReleaseDate   *time.Time `json:"release_date"`
}
type CreateTrackResponse struct {
	AlbumID string `json:"albumID"`
}

type GetAlbumRequest struct {
	AlbumID string `json:"albumID"`
}

type GetAlbumResponse struct {
	ID            string     `json:"id"`
	Title         string     `json:"title"`
	CoverImageURL *string    `json:"cover_image_url"`
	ReleaseDate   *time.Time `json:"release_date"`
	CreatedAt     time.Time  `json:"created_at"`
}

type Handler struct {
	createAlbumService in.CreateAlbumService
	getAlbumService    in.GetAlbumService
	listAlbumService   in.ListAlbumsService
}

func NewHandler(createAlbumService in.CreateAlbumService, getAlbumService in.GetAlbumService, listAlbumSerivce in.ListAlbumsService) *Handler {
	return &Handler{
		createAlbumService: createAlbumService,
		getAlbumService:    getAlbumService,
		listAlbumService:   listAlbumSerivce,
	}
}

func (h *Handler) CreateAlbum(w http.ResponseWriter, r *http.Request) {
	var req CreateAlbumRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}
	input := in.CreateAlbumInput{
		Title:         req.Title,
		CoverImageURL: req.CoverImageURL,
		ReleaseDate:   req.ReleaseDate,
	}
	out, err := h.createAlbumService.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp := CreateTrackResponse{
		AlbumID: out.AlbumID,
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) GetAlbumByID(w http.ResponseWriter, r *http.Request, albumID string) {

	out, err := h.getAlbumService.Execute(r.Context(), in.GetAlbumInput{AlbumID: albumID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := GetAlbumResponse{
		ID:            out.ID,
		Title:         out.Title,
		CoverImageURL: out.CoverImageURL,
		ReleaseDate:   out.ReleaseDate,
		CreatedAt:     out.CreatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) ListAlbums(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	input := in.ListAlbumsInput{
		Title: query.Get("title"),
	}

	out, err := h.listAlbumService.Execute(r.Context(), input)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	resp := []GetAlbumResponse{}

	for _, each := range out {
		curr := GetAlbumResponse{
			ID:            each.ID,
			Title:         each.Title,
			CoverImageURL: each.CoverImageURL,
			ReleaseDate:   each.ReleaseDate,
			CreatedAt:     each.CreatedAt,
		}
		resp = append(resp, curr)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
