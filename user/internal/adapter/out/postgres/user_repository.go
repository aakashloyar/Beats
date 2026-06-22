package postgres

import (
	"database/sql"

	"github.com/aakashloyar/beats/user/internal/application/ports/out"
	"github.com/aakashloyar/beats/user/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) out.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Save(user domain.User) error {
	query := `
		INSERT INTO users (
			id,
			username,
			email,
			created_at
		)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(query, user.ID, user.Username, user.Email, user.CreatedAt)
	return err
}

func (r *UserRepository) FindByID(userID string) (domain.User, error) {
	query := `
		SELECT
			id,
			username,
			email,
			created_at
		FROM users
		WHERE id = $1
	`

	row := r.db.QueryRow(query, userID)

	var user domain.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt); err != nil {
		return domain.User{}, err
	}

	return user, nil
}
