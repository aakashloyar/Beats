package http

import (
	"encoding/json"
	"github.com/aakashloyar/beats/track/internal/application/ports/in/track"
	"github.com/aakashloyar/beats/track/internal/domain"
	"net/http"
	"time"
)

type CreateTrackRequest struct {
	Title         string          `json:"title"`
	ArtistID      string          `json:"artist_id"`
	AlbumID       *string         `json:"album_id,omitempty"`
	CoverImageURL *string         `json:"cover_image_url,omitempty"`
	DurationMS    int64           `json:"duration_ms"`
	Language      domain.Language `json:"language"`
	ReleaseDate   *time.Time      `json:"release_date,omitempty"`
}

type CreateTrackResponse struct {
	TrackID string `json:"track_id"`
}

type GetTrackResponse struct {
	ID            string          `json:"id"`
	Title         string          `json:"title"`
	ArtistID      string          `json:"artist_id"`
	AlbumID       *string         `json:"album_id,omitempty"`
	CoverImageURL *string         `json:"cover_image_url,omitempty"`
	DurationMS    int64           `json:"duration_ms"`
	Language      domain.Language `json:"language"`
	ReleaseDate   *time.Time      `json:"release_data,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
}

type Handler struct {
	createTrackService in.CreateTrackService
	getTrackService    in.GetTrackService
	listTracksService  in.ListTracksService
}

func NewHandler(createTrackService in.CreateTrackService, getTrackService in.GetTrackService, listTracksService in.ListTracksService) *Handler {
	return &Handler{
		createTrackService: createTrackService,
		getTrackService:    getTrackService,
		listTracksService:  listTracksService,
	}
}

func (h *Handler) CreateTrack(w http.ResponseWriter, r *http.Request) {
	var req CreateTrackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
	}
	input := in.CreateTrackInput{
		Title:         req.Title,
		ArtistID:      req.ArtistID,
		AlbumID:       req.AlbumID,
		CoverImageURL: req.CoverImageURL,
		DurationMS:    req.DurationMS,
		Language:      req.Language,
		ReleaseDate:   req.ReleaseDate,
	}
	out, err := h.createTrackService.Execute(r.Context(), input)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := CreateTrackResponse{
		TrackID: out.TrackID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) GetTrackByID(w http.ResponseWriter, r *http.Request, trackID string) {
	out, err := h.getTrackService.Execute(r.Context(), trackID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := GetTrackResponse{
		ID:            out.ID,
		Title:         out.Title,
		ArtistID:      out.ArtistID,
		AlbumID:       out.AlbumID,
		CoverImageURL: out.CoverImageURL,
		DurationMS:    out.DurationMS,
		Language:      out.Language,
		ReleaseDate:   out.ReleaseDate,
		CreatedAt:     out.CreatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) ListTracks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	input := in.ListTracksInput{
		Title:    query.Get("title"),
		ArtistID: query.Get("artist_id"),
		AlbumID:  query.Get("album_id"),
		Limit:    query.Get("limit"),
		Offset:   query.Get("offset"),
	}

	out, err := h.listTracksService.Execute(r.Context(), input)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	resp := []GetTrackResponse{}

	for _, each := range out {
		curr := GetTrackResponse{
			ID:            each.ID,
			Title:         each.Title,
			ArtistID:      each.ArtistID,
			AlbumID:       each.AlbumID,
			CoverImageURL: each.CoverImageURL,
			DurationMS:    each.DurationMS,
			Language:      each.Language,
			ReleaseDate:   each.ReleaseDate,
			CreatedAt:     each.CreatedAt,
		}
		resp = append(resp, curr)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
