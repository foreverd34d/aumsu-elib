package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password *string
	DBName   string
	SSLMode  string
}

func NewDB(ctx context.Context, cfg Config) (*sqlx.DB, error) {
	connectString := fmt.Sprintf("host=%v port=%v user=%v dbname=%v sslmode=%v",
		cfg.Host, cfg.Port, cfg.User, cfg.DBName, cfg.SSLMode)
	if cfg.Password != nil {
		connectString += fmt.Sprintf("password=%v", *cfg.Password)
	}
	return sqlx.ConnectContext(ctx, "postgres", connectString)
}
