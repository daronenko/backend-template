package infra

import (
	"fmt"

	"github.com/daronenko/backend-template/internal/app/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func NewPostgres(conf *config.Config) (*sqlx.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s password=%s",
		conf.Postgres.Host,
		conf.Postgres.Port,
		conf.Postgres.User,
		conf.Postgres.Database,
		conf.Postgres.SSLMode,
		conf.Postgres.Password,
	)

	db, err := sqlx.Connect(conf.Postgres.Driver, dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(conf.Postgres.MaxOpenConns)
	db.SetConnMaxLifetime(conf.Postgres.ConnMaxLifetime)
	db.SetMaxIdleConns(conf.Postgres.MaxIdleConns)
	db.SetConnMaxIdleTime(conf.Postgres.ConnMaxIdleTime)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
