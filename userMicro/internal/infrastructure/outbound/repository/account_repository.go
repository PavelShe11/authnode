package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/PavelShe11/studbridge/user/internal/entity"
	"github.com/PavelShe11/studbridge/user/internal/port"

	"github.com/jmoiron/sqlx"
)

type accountRepository struct {
	db *sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) port.AccountRepository {
	return &accountRepository{
		db: db,
	}
}

var _ port.AccountRepository = (*accountRepository)(nil)

func (a *accountRepository) CreateAccount(ctx context.Context, account entity.Account) error {
	query := "INSERT INTO account (first_name, last_name, email) VALUES (:first_name, :last_name, :email)"
	_, err := a.db.NamedExecContext(ctx, query, account)
	if err != nil {
		return err
	}
	return nil
}

func (a *accountRepository) GetAccountByEmail(ctx context.Context, email string) (*entity.Account, error) {
	account := entity.Account{}
	query := "SELECT * FROM account WHERE email=$1"
	row := a.db.QueryRowxContext(ctx, query, email)
	err := row.StructScan(&account)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (a *accountRepository) GetAccountById(ctx context.Context, id string) (*entity.Account, error) {
	account := entity.Account{}
	query := "SELECT * FROM account WHERE id=$1"
	row := a.db.QueryRowxContext(ctx, query, id)
	err := row.StructScan(&account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}
