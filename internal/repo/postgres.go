package repo

import (
	"fmt"
	"la-warehouse/config"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type Pg struct {
	*sqlx.DB
}

func NewPostgres(cfg config.Postgres) (*Pg, error) {
	connectionURL := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)
	open, err := sqlx.Open("pgx", connectionURL)
	if err != nil {
		return nil, err
	}

	open.SetMaxOpenConns(cfg.Settings.MaxOpenConns)
	open.SetConnMaxLifetime(cfg.Settings.ConnMaxLifeTime * time.Second)
	open.SetMaxIdleConns(cfg.Settings.MaxIdleConns)
	open.SetConnMaxIdleTime(cfg.Settings.MaxIdleLifeTime * time.Second)

	return &Pg{open}, nil
}
