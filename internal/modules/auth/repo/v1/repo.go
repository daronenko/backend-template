package repo

import (
	"context"

	"github.com/daronenko/backend-template/internal/models"
	"github.com/daronenko/backend-template/pkg/utils"
	"github.com/google/uuid"
)

type Repo interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, userID uuid.UUID) error

	GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, user *models.User) (*models.User, error)

	FindByName(ctx context.Context, name string, query *utils.PaginationQuery) (*models.UsersList, error)
	GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*models.UsersList, error)
}

type Cache interface {
	Set(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	Delete(ctx context.Context, userID uuid.UUID) error
}
