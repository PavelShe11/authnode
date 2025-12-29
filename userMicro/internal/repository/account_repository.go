package repository

import (
	"database/sql"
	"errors"

	"github.com/PavelShe11/studbridge/user/internal/entity"

	"github.com/jmoiron/sqlx"
)

type AccountRepository struct {
	db *sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (a *AccountRepository) CreateAccount(account entity.Account) error {
	query := "INSERT INTO account (first_name, last_name, email) VALUES (:first_name, :last_name, :email)"
	_, err := a.db.NamedExec(query, account)
	if err != nil {
		return err
	}
	return nil
}

func (a *AccountRepository) GetAccountByEmail(email string) (*entity.Account, error) {
	account := entity.Account{}
	query := "SELECT * FROM account WHERE email=$1"
	row := a.db.QueryRowx(query, email)
	err := row.StructScan(&account)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (a *AccountRepository) GetAccountById(id string) (*entity.Account, error) {
	account := entity.Account{}
	query := "SELECT * FROM account WHERE id=$1"
	row := a.db.QueryRowx(query, id)
	err := row.StructScan(&account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}
