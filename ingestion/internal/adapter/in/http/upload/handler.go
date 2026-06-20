package http

import (
	"encoding/json"
	"github.com/aakashloyar/beats/ingestion/internal/application/ports/in/upload"
	"net/http"
	"time"
	"github.com/aakashloyar/beats/config"
)

type Handler struct {
	initUploadService in.InitUploadService
	completeUploadService in.CompleteUploadService
}

type initUploadRequest struct {
	Title         string          `json:"title"`
	ArtistIDs     []string        `json:"artist_ids"`
	AlbumID       *string         `json:"album_id,omitempty"`
	CoverImageURL *string         `json:"cover_image_url,omitempty"`
	Language      config.Language `json:"language"`
	ReleasedAt    *time.Time      `json:"released_at,omitempty"`
	FileName      string          `json:"file_name"`
	FileSize      int64           `json:"file_size"`
}

type initUploadResponse struct {
	UploadID     string      `json:"upload_id"`
	MaxChunkSize int64       `json:"max_chunk_size"`
	UploadURLs   []UploadURL `json:"chunks"`
}
type UploadURL struct {
	ChunkNumber int   `json:"chunk_number"`
	URL        string `json:"url"`
}

type completeUploadRequest struct {
	UploadID string `json:"upload_id"`
}
type completUploadResponse struct {

}

func NewHandler(initUploadService in.InitUploadService, completeUploadService in.CompleteUploadService) *Handler {
	return &Handler{
		initUploadService: initUploadService,
	}
}

func (h *Handler) InitUpload(w http.ResponseWriter, r *http.Request) {
	var req initUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
	}
	input := in.InitUploadInput{
		Title:         req.Title,
		ArtistIDs:     req.ArtistIDs,
		AlbumID:       req.AlbumID,
		CoverImageURL: req.CoverImageURL,
		ReleasedAt:    req.ReleasedAt,
		FileName:      req.FileName,
		FileSize:      req.FileSize,
	}

	out, err := h.initUploadService.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uploadURLs := make([]UploadURL,0,len(out.UploadURLs))
	for _,u := range out.UploadURLs {
		uploadURLs = append(uploadURLs, UploadURL{
			ChunkNumber: u.ChunkNumber,
			URL:         u.URL, 
		})
	}
	res := initUploadResponse{
		UploadID: out.UploadID,
		MaxChunkSize: out.MaxChunkSize,
		UploadURLs: uploadURLs,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}


func (h *Handler) CompleteUpload(w http.ResponseWriter, r *http.Request) {
	var req completeUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
	}
	input := in.CompleteUploadInput{
		UploadID: req.UploadID,
	}

	_, err := h.completeUploadService.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := completUploadResponse{}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}
