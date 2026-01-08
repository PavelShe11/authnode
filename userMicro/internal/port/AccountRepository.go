package port

import (
	"context"

	"github.com/PavelShe11/studbridge/user/internal/entity"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, account entity.Account) error
	GetAccountByEmail(ctx context.Context, email string) (*entity.Account, error)
	GetAccountById(ctx context.Context, id string) (*entity.Account, error)
}
