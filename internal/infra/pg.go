package infra

import (
	"fmt"
	"time"

	"github.com/daronenko/backend-template/internal/app/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	maxOpenConns    = 60
	connMaxLifetime = 120
	maxIdleConns    = 30
	connMaxIdleTime = 20
)

func NewPostgres(cfg *config.Config) (*sqlx.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s password=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Database,
		cfg.Postgres.SSLMode,
		cfg.Postgres.Password,
	)

	db, err := sqlx.Connect(cfg.Postgres.Driver, dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifetime * time.Second)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(connMaxIdleTime * time.Second)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
