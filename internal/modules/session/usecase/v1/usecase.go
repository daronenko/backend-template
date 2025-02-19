package usecase

import (
	"context"

	"github.com/daronenko/backend-template/internal/models"
	"github.com/google/uuid"
)

type Usecase interface {
	Create(ctx context.Context, session *models.Session) (string, error)
	GetByID(ctx context.Context, sessionID uuid.UUID) (*models.Session, error)
	DeleteByID(ctx context.Context, sessionID uuid.UUID) error
}
