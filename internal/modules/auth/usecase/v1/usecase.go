package usecase

import (
	"context"

	"github.com/daronenko/backend-template/internal/models"
	"github.com/google/uuid"
)

type Usecase interface {
	Register(ctx context.Context, user *models.User) (*models.UserWithToken, error)
	Login(ctx context.Context, user *models.User) (*models.UserWithToken, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, userID uuid.UUID) error
	GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
}
