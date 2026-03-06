package port

import (
	"context"

	"github.com/PavelShe11/studbridge/authMicro/internal/entity"
)

// AccountProvider - интерфейс для взаимодействия с Account сервисом (userMicro)
// Service layer будет зависеть от этого интерфейса вместо конкретного gRPC клиента
type AccountProvider interface {
	// ValidateAccountData проверяет корректность данных аккаунта через внешний сервис
	// userData - динамические данные (map[string]interface{})
	// lang - язык для локализации ошибок
	ValidateAccountData(ctx context.Context, userData map[string]interface{}, lang string) error

	// CreateAccount создает новый аккаунт через внешний сервис
	CreateAccount(ctx context.Context, userData map[string]interface{}, lang string) error

	// GetAccountByEmail возвращает данные аккаунта (или nil если не найден)
	GetAccountByEmail(ctx context.Context, email string) (*entity.Account, error)

	// GetAccessTokenPayload возвращает map с данными для JWT access token
	// Возвращает: map с claims для токена
	GetAccessTokenPayload(ctx context.Context, accountId string) (map[string]interface{}, error)
}
