package repository

import (
	"context"
	"database/sql"

	"github.com/PavelShe11/studbridge/auth/internal/domain"
	"github.com/jmoiron/sqlx"
)

type RefreshTokenSessionRepository struct {
	db *sqlx.DB
}

func NewRefreshTokenSessionRepository(db *sqlx.DB) *RefreshTokenSessionRepository {
	return &RefreshTokenSessionRepository{db: db}
}

// Save сохраняет новую сессию refresh token
func (r *RefreshTokenSessionRepository) Save(session *domain.RefreshTokenSession) error {
	query := `
		INSERT INTO refresh_token_session (account_id, refresh_token, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	return r.db.QueryRowx(query,
		session.AccountID,
		session.RefreshToken,
		session.ExpiresAt,
	).Scan(&session.Id, &session.CreatedAt)
}

// FindByToken находит сессию по токену
func (r *RefreshTokenSessionRepository) FindByToken(token string) (*domain.RefreshTokenSession, error) {
	var session domain.RefreshTokenSession
	query := `SELECT * FROM refresh_token_session WHERE refresh_token = $1`
	err := r.db.Get(&session, query, token)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &session, err
}

// DeleteByToken удаляет сессию по токену
func (r *RefreshTokenSessionRepository) DeleteByToken(token string) error {
	query := `DELETE FROM refresh_token_session WHERE refresh_token = $1`
	_, err := r.db.Exec(query, token)
	return err
}

// CleanExpired удаляет истекшие сессии
func (r *RefreshTokenSessionRepository) CleanExpired(ctx context.Context) error {
	query := `DELETE FROM refresh_token_session WHERE expires_at < NOW()`
	_, err := r.db.ExecContext(ctx, query)
	return err
}
