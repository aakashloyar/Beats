package postgres

import (
	"database/sql"

	"github.com/aakashloyar/beats/track/internal/application/ports/out"
	"github.com/aakashloyar/beats/track/internal/domain"
)

type TrackRepository struct {
	db *sql.DB 
}

func NewTrackRepository(db *sql.DB) out.TrackRepository {
	return &TrackRepository{db: db}
}

func (r *TrackRepository) Save(track *domain.Track) error {

	query := `
		INSERT INTO tracks (
		    id,
			title,
			artist_id,
			album_id,
			cover_image_url,
			duration_ms,
			language,
			release_data,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)	
	`
	_, err := r.db.Exec(
		query, 
		track.ID,
		track.Title,
		track.ArtistID,
		track.AlbumID,
		track.CoverImageURL,
		track.DurationMS,
		track.Language,
		track.ReleaseDate,
		track.CreatedAt,
	)
	return err 
}

func (r *TrackRepository) FindById(id string) (*domain.Track, error) {
	query := `
		SELECT 
			id,
			title,
			artist_id,
			album_id,
			coverimage_url,
			duration_ms,
			language,
			release_date,
			created_at
		FROM Tracks
		WHERE id = $1	
	`
	row := r.db.QueryRow(query, id)

	var track domain.Track
	err := row.Scan(
		&track.ID,
		&track.Title,
		&track.ArtistID,
		&track.AlbumID,
		&track.CoverImageURL,
		&track.DurationMS,
		&track.Language,
		&track.ReleaseDate,
		&track.CreatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &track, nil 
}


func (r *TrackRepository) FindByArtist(artistID string) ([]*domain.Track, error) {
	query := `
		SELECT
			id,
			title,
			artist_id,
			album_id,
			cover_image_url,
			duration_ms,
			language,
			release_date,
			created_at
		FROM tracks
		WHERE artist_id = $1
		ORDER BY release_date DESC
	`

	rows, err := r.db.Query(query, artistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []*domain.Track

	for rows.Next() {
		var track domain.Track

		err := rows.Scan(
			&track.ID,
			&track.Title,
			&track.ArtistID,
			&track.AlbumID,
			&track.CoverImageURL,
			&track.DurationMS,
			&track.Language,
			&track.ReleaseDate,
			&track.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		tracks = append(tracks, &track)
	}

	return tracks, nil
}
