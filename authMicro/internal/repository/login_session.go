package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/PavelShe11/studbridge/auth/internal/domain"

	"github.com/jmoiron/sqlx"
)

type LoginSessionRepository struct {
	db *sqlx.DB
}

func NewLoginSessionRepository(db *sqlx.DB) *LoginSessionRepository {
	return &LoginSessionRepository{
		db: db,
	}
}

func (r *LoginSessionRepository) FindByEmail(email string) (*domain.LoginSession, error) {
	query := "SELECT * FROM login_session WHERE email = $1"
	result := &domain.LoginSession{}
	row := r.db.QueryRowx(query, email)
	err := row.StructScan(result)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *LoginSessionRepository) Save(session *domain.LoginSession) error {
	query := `INSERT INTO login_session (account_id, email, code, code_expires)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (email)
	DO UPDATE
	SET account_id = EXCLUDED.account_id, code = EXCLUDED.code, code_expires = EXCLUDED.code_expires
	RETURNING id, account_id, email, code, code_expires, created_at`

	err := r.db.QueryRowx(query, session.AccountId, session.Email, session.Code, session.CodeExpires).StructScan(session)
	if err != nil {
		return err
	}
	return nil
}

func (r *LoginSessionRepository) DeleteByEmail(context context.Context, email string) error {
	query := "DELETE FROM login_session WHERE email = $1"
	_, err := r.db.ExecContext(context, query, email)
	if err != nil {
		return err
	}
	return nil
}

func (r *LoginSessionRepository) CleanExpired(ctx context.Context) error {
	query := "DELETE FROM login_session WHERE code_expires < NOW()"
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	return nil
}
