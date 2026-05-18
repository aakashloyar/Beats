package postgres

import (
	"database/sql"
	"github.com/aakashloyar/beats/track/internal/application/ports/out"
	"github.com/aakashloyar/beats/track/internal/domain"
)

type AlbumRepository struct {
	db *sql.DB
}

func NewAlbumRepository(db *sql.DB) out.AlbumRepository {
	return &AlbumRepository{db: db}
}

func (r *AlbumRepository) Save(input domain.Album) error {
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
		input.Title,
		input.CoverImageURL,
		input.CreatedAt,
	)
	return err
}

func (r *AlbumRepository) FindByID(albumID string) (domain.Album, error) {
	query := `
		SELECT
		    id,
			title, 
			cover_image_url,
			release_date,
			created_at
		FROM albums
		WHERE id = $1
	`
	row := r.db.QueryRow(query, albumID)

	var album domain.Album
	err := row.Scan(
		&album.ID,
		&album.Title,
		&album.CoverImageURL,
		&album.ReleaseDate,
		&album.CreatedAt,
	)

	if err != nil {
		return domain.Album{}, err
	}

	return album, nil
}

func (r *AlbumRepository) ListAlbums(title string) ([]domain.Album, error) {
	query := `
	    SELECT 
		    id, 
			title, 
			cover_image_url, 
			release_date, 
			created_at
	    FROM albums
		WHERE LOWER(title) LIKE LOWER($1)
	`

	rows, err := r.db.Query(query, title)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var albums []domain.Album

	for rows.Next() {
		var album domain.Album

		err := rows.Scan(
			&album.ID,
			&album.Title,
			&album.CoverImageURL,
			&album.ReleaseDate,
			&album.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		albums = append(albums, album)
	}

	return albums, nil
}
