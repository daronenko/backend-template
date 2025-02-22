package repo

import (
	"context"

	"github.com/daronenko/backend-template/internal/model/v1"
	"github.com/google/uuid"
)

type Repo interface {
	Create(ctx context.Context, session *model.Session) (string, error)
	GetByID(ctx context.Context, sessionID uuid.UUID) (*model.Session, error)
	DeleteByID(ctx context.Context, sessionID uuid.UUID) error
}
