package usecase

import "context"

type HealthUsecase interface {
	Ping(ctx context.Context) error
}
