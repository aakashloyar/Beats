package http

import (
	"encoding/json"
	"github.com/aakashloyar/beats/track/internal/application/ports/in/track"
	"github.com/aakashloyar/beats/config"
	"net/http"
	"time"
)

type GetTrackResponse struct {
	ID            string            `json:"id"`
	Title         string            `json:"title"`
	ArtistIDs     []string          `json:"artist_ids"`
	AlbumID       *string           `json:"album_id,omitempty"`
	CoverImageURL *string           `json:"cover_image_url,omitempty"`
	DurationMS    int64             `json:"duration_ms"`
	Language      []config.Language `json:"language"`
	ReleasedAt    *time.Time        `json:"released_at,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
}

type ListAudioVariantsByTrackResponse struct {
	ID           string    `json:"id"`
	TrackID      string    `json:"track_id"`
	Codec        string    `json:"codec"`
	BitrateKbps  int       `json:"bitrate_kbps"`
	SampleRateHz int       `json:"sample_rate_hz"`
	Channels     int       `json:"channels"`
	DurationMs   int64     `json:"duration_ms"`
	FileURL      string    `json:"file_url"`
	CreatedAt    time.Time `json:"created_at"`
}

type Handler struct {
	createTrackService      in.CreateTrackService
	getTrackService         in.GetTrackService
	listTracksService       in.ListTracksService
}

func NewHandler(createTrackService in.CreateTrackService, getTrackService in.GetTrackService, listTracksService in.ListTracksService) *Handler {
	return &Handler{
		createTrackService:      createTrackService,
		getTrackService:         getTrackService,
		listTracksService:       listTracksService,
	}
}

func (h *Handler) GetTrackByID(w http.ResponseWriter, r *http.Request, trackID string) {
	out, err := h.getTrackService.Execute(r.Context(), in.GetTrackInput{TrackID: trackID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := GetTrackResponse{
		ID:            out.ID,
		Title:         out.Title,
		ArtistIDs:     out.ArtistIDs,
		AlbumID:       out.AlbumID,
		CoverImageURL: out.CoverImageURL,
		DurationMS:    out.DurationMS,
		Language:      out.Language,
		ReleasedAt:    out.ReleasedAt,
		CreatedAt:     out.CreatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) ListTracks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	input := in.ListTracksInput{
		Title:     query.Get("title"),
		ArtistIDs: []string{query.Get("artist_id")},
		AlbumID:   query.Get("album_id"),
		Limit:     query.Get("limit"),
		Offset:    query.Get("offset"),
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
			ArtistIDs:     each.ArtistIDs,
			AlbumID:       each.AlbumID,
			CoverImageURL: each.CoverImageURL,
			DurationMS:    each.DurationMS,
			Language:      each.Language,
			ReleasedAt:    each.ReleasedAt,
			CreatedAt:     each.CreatedAt,
		}
		resp = append(resp, curr)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}