package domain

import "time"

type RefreshTokenSession struct {
	Id           string    `db:"id"`
	AccountID    string    `db:"account_id"`
	RefreshToken string    `db:"refresh_token"`
	ExpiresAt    time.Time `db:"expires_at"`
	CreatedAt    time.Time `db:"created_at"`
}
