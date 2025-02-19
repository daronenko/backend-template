package models

import "github.com/google/uuid"

type Session struct {
	ID     uuid.UUID `json:"id" redis:"id"`
	UserID uuid.UUID `json:"user_id" redis:"user_id"`
}
