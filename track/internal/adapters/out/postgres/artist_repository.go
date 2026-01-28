package postgres

import (
	"database/sql"
	"github.com/aakashloyar/beats/track/internal/application/ports/out"
	"github.com/aakashloyar/beats/track/internal/domain"
)

type ArtistRepository struct {
	db *sql.DB
}

func NewArtistRepository(db *sql.DB) out.ArtistRepository {
	return &ArtistRepository{db: db}
}

func (r *ArtistRepository) Save(input domain.Artist) error {
	query := `
	    INSERT INTO artists (
		    id,
			name,
			bio,
			profile_image_url,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(
		query,
		input.ID,
		input.Name,
		input.Bio,
		input.CreatedAt,
	)
	return err
}

func (r *ArtistRepository) FindByID(artistID string) (domain.Artist, error) {
	query := `
		SELECT
		    id,
			name, 
			bio,
			profile_image_url,
			created_at
		FROM artists 
		WHERE id = $1
	`
	row := r.db.QueryRow(query, artistID)

	var artist domain.Artist
	err := row.Scan(
		&artist.ID,
		&artist.Name,
		&artist.Bio,
		&artist.ProfileImageURL,
		&artist.CreatedAt,
	)

	if err != nil {
		return domain.Artist{}, err
	}

	return artist, nil
}
