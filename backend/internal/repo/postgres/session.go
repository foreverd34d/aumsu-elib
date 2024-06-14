package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/foreverd34d/aumsu-elib/internal/errs"
	"github.com/foreverd34d/aumsu-elib/internal/model"

	"github.com/jmoiron/sqlx"
)

type TokenPsqlRepo struct {
	db *sqlx.DB
}

func NewTokenPsqlRepo(db *sqlx.DB) *TokenPsqlRepo {
	return &TokenPsqlRepo{db: db}
}

func (tr *TokenPsqlRepo) Create(ctx context.Context, userID int, input *model.NewToken) (*model.Token, error) {
	txCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	tx, err := tr.db.BeginTxx(txCtx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin the transaction: %w: %w", errs.Internal, err)
	}

	var sessionID int
	sessionQuery := `
		INSERT INTO sessions (user_id)
		VALUES ($1)
		RETURNING session_id
	`
	if err = tx.GetContext(ctx, &sessionID, sessionQuery, userID); err != nil {
		return nil, fmt.Errorf("insert the session: %w: %w", errs.Internal, err)
	}

	token := new(model.Token)
	tokenQuery := `
		INSERT INTO tokens (refresh_token, expires_at, session_id)
		VALUES ($1, $2, $3)
		RETURNING token_id, refresh_token, expires_at, session_id
	`
	if err := tx.GetContext(ctx, token, tokenQuery, input.RefreshToken, input.ExpiresAt, sessionID); err != nil {
		return nil, fmt.Errorf("insert the token: %w: %w", errs.Internal, err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit the changes: %w: %w", errs.Internal, err)
	}
	return token, nil
}

func (tr *TokenPsqlRepo) PopByRefreshToken(ctx context.Context, refreshToken string) (*model.Token, error) {
	token := new(model.Token)
	tokenQuery := `
		DELETE FROM tokens
		WHERE refresh_token = $1
		RETURNING token_id, refresh_token, expires_at, session_id
	`
	if err := tr.db.GetContext(ctx, token, tokenQuery, refreshToken); err != nil {
		return nil, fmt.Errorf("pop the refresh token: %w: %w", errs.Internal, err)
	}
	return token, nil
}

func (tr *TokenPsqlRepo) UpdateRefreshToken(ctx context.Context, sessionID int, update *model.NewToken) (*model.Token, error) {
	token := new(model.Token)
	query := `
		INSERT INTO tokens (refresh_token, expires_at, session_id)
		VALUES ($1, $2, $3)
		RETURNING token_id, refresh_token, expires_at, session_id
	`
	if err := tr.db.GetContext(ctx, token, query, update.RefreshToken, update.ExpiresAt, sessionID); err != nil {
		return nil, fmt.Errorf("update the refresh token: %w: %w", errs.Internal, err)
	}
	return token, nil
}

func (tr *TokenPsqlRepo) EndSession(ctx context.Context, sessionID int) error {
	query := `
		UPDATE sessions
		SET logged_out_at = $1
		WHERE session_id = $2
	`
	_, err := tr.db.ExecContext(ctx, query, time.Now(), sessionID)
	if err != nil {
		return fmt.Errorf("end the session: %w: %w", errs.Internal, err)
	}
	return nil
}
func (tr *TokenPsqlRepo) GetUserFromSession(ctx context.Context, sessionID int) (*model.User, error) {
	user := new(model.User)
	query := `
		SELECT u.*
		FROM users u
		JOIN sessions s USING(user_id)
		WHERE s.session_id = $1
	`
	if err := tr.db.GetContext(ctx, user, query, sessionID); err != nil {
		baseErr := errs.Internal
		if errors.Is(err, sql.ErrNoRows) {
			baseErr = errs.NotFound
		}
		return nil, fmt.Errorf("get the user: %w: %w", baseErr, err)
	}
	return user, nil
}
