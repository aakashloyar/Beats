package in

import (
	"context"
	"time"
)

type GetUserInput struct {
	UserID string
}

type GetUserOutput struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type GetUserService interface {
	Execute(ctx context.Context, input GetUserInput) (GetUserOutput, error)
}
