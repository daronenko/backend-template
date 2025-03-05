package usecase

import (
	"context"

	"github.com/daronenko/backend-template/internal/model/v1"
	"github.com/daronenko/backend-template/internal/util"
	"github.com/google/uuid"
)

type Usecase interface {
	Register(ctx context.Context, user *model.User) (*model.UserWithToken, error)
	Login(ctx context.Context, user *model.User) (*model.UserWithToken, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
	Delete(ctx context.Context, userID uuid.UUID) error
	GetByID(ctx context.Context, userID uuid.UUID) (*model.User, error)
	FindByUsername(ctx context.Context, name string, query *util.PaginationQuery) (*model.UsersList, error)
	GetUsers(ctx context.Context, pq *util.PaginationQuery) (*model.UsersList, error)
}
