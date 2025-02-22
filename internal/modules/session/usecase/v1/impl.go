package usecase

import (
	"context"

	"github.com/daronenko/backend-template/internal/model/v1"
	"github.com/daronenko/backend-template/internal/modules/session/repo/v1"
	"github.com/google/uuid"
)

// User session usecase
type Session struct {
	repo repo.Repo
}

func NewSession(sessionRepo repo.Repo) Usecase {
	return &Session{sessionRepo}
}

// Create new user session
func (s *Session) Create(ctx context.Context, session *model.Session) (string, error) {
	return s.repo.Create(ctx, session)
}

// Get user session by session id
func (s *Session) GetByID(ctx context.Context, sessionID uuid.UUID) (*model.Session, error) {
	return s.repo.GetByID(ctx, sessionID)
}

// Delete user session by session id
func (s *Session) DeleteByID(ctx context.Context, sessionID uuid.UUID) error {
	return s.repo.DeleteByID(ctx, sessionID)
}
