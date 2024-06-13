package postgres

import (
	"context"
	"time"

	"github.com/foreverd34d/aumsu-elib/internal/model"

	"github.com/jmoiron/sqlx"
)

type TokenPostgresRepo struct {
	db *sqlx.DB
}

func NewTokenPostgesRepo(db *sqlx.DB) *TokenPostgresRepo {
	return &TokenPostgresRepo{db: db}
}

func (tr *TokenPostgresRepo) Create(ctx context.Context, userID int, input *model.NewToken) (*model.Token, error) {
	tx, err := tr.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	var sessionID int
	sessionQuery := `
		INSERT INTO sessions (user_id)
		VALUES ($1)
		RETURNING session_id
	`
	if err = tx.GetContext(ctx, &sessionID, sessionQuery, userID); err != nil {
		tx.Rollback()
		return nil, err
	}

	token := new(model.Token)
	tokenQuery := `
		INSERT INTO tokens (refresh_token, expires_at, session_id)
		VALUES ($1, $2, $3)
		RETURNING token_id, refresh_token, expires_at, session_id
	`
	if err := tx.GetContext(ctx, token, tokenQuery, input.RefreshToken, token.ExpiresAt, sessionID); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return token, nil
}

func (tr *TokenPostgresRepo) PopByRefreshToken(ctx context.Context, refreshToken string) (*model.Token, error) {
	token := new(model.Token)
	tokenQuery := `
		DELETE FROM tokens
		WHERE refresh_token = $1
		RETURNING token_id, refresh_token, expires_at, session_id
	`
	err := tr.db.GetContext(ctx, token, tokenQuery, refreshToken)
	return token, err
}

func (tr *TokenPostgresRepo) UpdateRefreshToken(ctx context.Context, sessionID int, update *model.NewToken) (*model.Token, error) {
	token := new(model.Token)
	query := `
		INSERT INTO tokens (refresh_token, expires_at, session_id)
		VALUES ($1, $2, $3)
		RETURNING token_id, refresh_token, expires_at, session_id
	`
	err := tr.db.GetContext(ctx, token, query, update.RefreshToken, update.ExpiresAt, sessionID)
	return token, err
}

func (tr *TokenPostgresRepo) EndSession(ctx context.Context, sessionID int) error {
	query := `
		UPDATE sessions
		SET logged_out_at = $1
		WHERE session_id = $2
	`
	_, err := tr.db.ExecContext(ctx, query, time.Now(), sessionID)
	return err
}
func (tr *TokenPostgresRepo) GetUserFromSession(ctx context.Context, sessionID int) (*model.User, error) {
	user := new(model.User)
	query := `
		SELECT u.*
		FROM users u
		JOIN sessions s USING(user_id)
		WHERE s.session_id = $1
	`
	err := tr.db.GetContext(ctx, user, query, sessionID)
	return user, err
}
// func (sr *TokenPostgresRepo) Create(ctx context.Context, input *model.NewToken) (*model.Token, error) {
// 	tx, err := sr.db.BeginTxx(ctx, nil)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	token := new(model.Token)
// 	tokenQuery := `
// 		INSERT INTO tokens (refresh_token, expires_at, user_id)
// 		VALUES ($1, $2, $3)
// 		RETURNING session_id, refresh_token, expires_at, user_id
// 	`
// 	if err := tx.GetContext(ctx, token, tokenQuery, input.RefreshToken, input.ExpiresAt, input.UserID); err != nil {
// 		tx.Rollback()
// 		return nil, err
// 	}
//
// 	sessionQuery := `INSERT INTO sessions (user_id) VALUES ($1)`
// 	_, err = tx.ExecContext(ctx, sessionQuery, input.UserID)
// 	if err != nil {
// 		tx.Rollback()
// 		return nil, err
// 	}
//
// 	if err := tx.Commit(); err != nil {
// 		return nil, err
// 	}
// 	return token, nil
// }
//
// // func (sr *TokenPostgresRepo) Update(ctx context.Context)
//
// func (sr *TokenPostgresRepo) PopByRefreshToken(ctx context.Context, refreshToken string) (*model.Token, error) {
// 	token := new(model.Token)
// 	tokenQuery := `
// 		DELETE FROM tokens
// 		WHERE refresh_token = $1
// 		RETURNING token_id, refresh_token, expires_at, user_id
// 	`
// 	err := sr.db.GetContext(ctx, token, tokenQuery, refreshToken)
// 	return token, err
// }
//
// func (sr *TokenPostgresRepo) DeleteToken(ctx context.Context, refreshToken string) error {
// 	tx, err := sr.db.BeginTxx(ctx, nil)
// 	if err != nil {
// 		return err
// 	}
//
// 	var userID int
// 	tokenQuery := `
// 		DELETE FROM tokens
// 		WHERE refresh_token = $1
// 		RETURNING user_id
// 	`
// 	if err := tx.GetContext(ctx, &userID, tokenQuery, refreshToken); err != nil {
// 		tx.Rollback()
// 		return err
// 	}
//
// 	sessionQuery := `
// 		UPDATE sessions
// 		SET logged_out_at = $1
// 		WHERE user_id = $2
// 	`
// 	_, err = tx.ExecContext(ctx, sessionQuery, time.Now(), userID)
// 	if err != nil {
// 		tx.Rollback()
// 		return err
// 	}
//
// 	if err := tx.Commit(); err != nil {
// 		return err
// 	}
// 	return nil
// }
