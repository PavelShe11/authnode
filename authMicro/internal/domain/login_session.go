package domain

import "time"

type LoginSession struct {
	Id          string    `db:"id"`
	AccountId   string    `db:"account_id"`
	Email       string    `db:"email"`
	Code        string    `db:"code"` // Stores bcrypt hash of verification code
	CodeExpires time.Time `db:"code_expires"`
	CreateAt    time.Time `db:"created_at"`
}
