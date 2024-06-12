package postgres

import (
	"context"

	"github.com/foreverd34d/aumsu-elib/internal/model"

	"github.com/jmoiron/sqlx"
)

type SessionPostgresRepo struct {
	db *sqlx.DB
}

func NewSessionPostgesRepo(db *sqlx.DB) *SessionPostgresRepo {
	return &SessionPostgresRepo{db: db}
}

func (sr *SessionPostgresRepo) Create(ctx context.Context, input *model.NewSession) (*model.Session, error) {
	session := new(model.Session)
	query := `
		INSERT INTO sessions (refresh_token, expires_at, user_id)
		VALUES ($1, $2, $3)
		RETURNING session_id, refresh_token, expires_at, user_id
	`
	err := sr.db.GetContext(ctx, session, query, input.RefreshToken, input.ExpiresAt, input.UserID)
	return session, err
}

func (sr *SessionPostgresRepo) PopByRefreshToken(ctx context.Context, refreshToken string) (*model.Session, error) {
	session := new(model.Session)
	query := `
		DELETE FROM sessions
		WHERE refresh_token = $1
		RETURNING session_id, refresh_token, expires_at, user_id
	`
	err := sr.db.GetContext(ctx, session, query, refreshToken)
	return session, err
}
