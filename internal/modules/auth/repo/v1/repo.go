package repo

import (
	"context"

	"github.com/daronenko/backend-template/internal/model/v1"
	"github.com/daronenko/backend-template/internal/util"
	"github.com/google/uuid"
)

type Repo interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
	Delete(ctx context.Context, userID uuid.UUID) error

	GetByID(ctx context.Context, userID uuid.UUID) (*model.User, error)
	GetByEmail(ctx context.Context, user *model.User) (*model.User, error)

	FindByName(ctx context.Context, name string, query *util.PaginationQuery) (*model.UsersList, error)
	GetUsers(ctx context.Context, pq *util.PaginationQuery) (*model.UsersList, error)
}

type Cache interface {
	Set(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, userID uuid.UUID) (*model.User, error)
	Delete(ctx context.Context, userID uuid.UUID) error
}
