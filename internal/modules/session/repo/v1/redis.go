package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"emperror.dev/errors"
	"github.com/daronenko/backend-template/internal/app/config"
	"github.com/daronenko/backend-template/internal/model/v1"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Repo for user sessions
type Session struct {
	redis *redis.Client
	conf  *config.Config
}

func NewSession(redis *redis.Client, conf *config.Config) Repo {
	return &Session{redis, conf}
}

// Create new user session
func (r *Session) Create(ctx context.Context, session *model.Session) (string, error) {
	session.ID = uuid.New()

	sessionBytes, err := json.Marshal(&session)
	if err != nil {
		return "", errors.Wrap(err, "repo.Session.Create.json.Marshal")
	}

	if err = r.redis.Set(ctx, r.cacheKey(session.ID), sessionBytes, time.Second*time.Duration(r.conf.App.Auth.Session.Cache.Expire)).Err(); err != nil {
		return "", errors.Wrap(err, "repo.Session.Set.redis.Set")
	}

	return session.ID.String(), nil
}

// Get user session by id
func (r *Session) GetByID(ctx context.Context, sessionID uuid.UUID) (*model.Session, error) {
	sessionStr, err := r.redis.Get(ctx, r.cacheKey(sessionID)).Result()
	if err != nil {
		return nil, errors.Wrap(err, "repo.Session.GetByID.redis.Get")
	}

	session := &model.Session{}
	if err = json.Unmarshal([]byte(sessionStr), session); err != nil {
		return nil, errors.Wrap(err, "repo.Session.GetByID.json.Unmarshal")
	}

	return session, nil
}

// Delete existing user session
func (r *Session) DeleteByID(ctx context.Context, sessionID uuid.UUID) error {
	if err := r.redis.Del(ctx, r.cacheKey(sessionID)).Err(); err != nil {
		return errors.Wrap(err, "repo.Session.DeleteByID.redis.Del")
	}

	return nil
}

func (r *Session) cacheKey(sessionID uuid.UUID) string {
	return fmt.Sprintf("%s:%s", r.conf.App.Auth.Session.Cache.Prefix, sessionID)
}
