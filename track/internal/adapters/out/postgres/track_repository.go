package postgres

import (
	"database/sql"

	"github.com/aakashloyar/beats/track/internal/application/ports/out"
	"github.com/aakashloyar/beats/track/internal/application/ports/in/track"
	"github.com/aakashloyar/beats/track/internal/domain"
	"strings"
)

type TrackRepository struct {
	db *sql.DB 
}

func NewTrackRepository(db *sql.DB) out.TrackRepository {
	return &TrackRepository{db: db}
}

func (r *TrackRepository) Save(track domain.Track) error {

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

func (r *TrackRepository) FindByID(id string) (domain.Track, error) {
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
		return domain.Track{}, err
	}
	return track, nil 
}


func (r *TrackRepository) ListTracks(input in.ListTracksInput) ([]domain.Track, error) {
	query := `
	SELECT id, title, artist_id, album_id
	FROM tracks
	`

	var (
		conditions []string
		args       []any
	)

	if input.Title != "" {
		conditions = append(conditions, "LOWER(title) LIKE LOWER(?)")
		args = append(args, "%"+input.Title+"%")
	}

	if input.ArtistID != "" {
		conditions = append(conditions, "artist_id = ?")
		args = append(args, input.ArtistID)
	}

	if input.AlbumID != "" {
		conditions = append(conditions, "album_id = ?")
		args = append(args, input.AlbumID)
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY created_at DESC "

	if input.Limit != "" {
		query += "LIMIT ?"
		args = append(args, input.Limit)
	}
	if input.Offset != "" {
		query += "OFFSET ?"
		args = append(args, query)
	}

	rows, err := r.db.Query(query, args)


	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tracks []domain.Track

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

		tracks = append(tracks, track)
	}

	return tracks, nil
}
