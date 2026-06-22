package out

import "github.com/aakashloyar/beats/user/internal/domain"

type UserRepository interface {
	Save(user domain.User) error
	FindByID(userID string) (domain.User, error)
}
