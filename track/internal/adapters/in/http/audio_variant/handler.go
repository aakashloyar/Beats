package http

import (
	"encoding/json"
	"github.com/aakashloyar/beats/track/internal/application/ports/in/audio_variant"
	"github.com/aakashloyar/beats/track/internal/domain"
	"net/http"
)

type CreateAudioVariantRequest struct {
	TrackID      string       `json:"track_id"`
	Codec        domain.Codec `json:"codec"`
	BitrateKbps  int          `json:"bitrate_kbps"`
	SampleRateHz int          `json:"sample_rate_hz"`
	Channels     int          `json:"channels"`
	DurationMs   int64        `json:"duration_ms"`
	FileURL      string       `json:"file_url"`
}

type CreateAudioVariantResponse struct {
	AudioVariantID string `json:"audio_variant_id"`
}

type Handler struct {
	createAudioVariant in.CreateAudioVariantService
}

func NewHandler(createAudioVariant in.CreateAudioVariantService) *Handler {
	return &Handler{
		createAudioVariant: createAudioVariant,
	}
}

func (h *Handler) CreateAudioVariant(w http.ResponseWriter, r *http.Request) {
	var req CreateAudioVariantRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	input := in.CreateAudioVariantInput{
		TrackID:      req.TrackID,
		Codec:        req.Codec,
		BitrateKbps:  req.BitrateKbps,
		SampleRateHz: req.SampleRateHz,
		Channels:     req.Channels,
		DurationMs:   req.DurationMs,
		FileURL:      req.FileURL,
	}

	out, err := h.createAudioVariant.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := CreateAudioVariantResponse{
		AudioVariantID: out.AudioVariantID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
