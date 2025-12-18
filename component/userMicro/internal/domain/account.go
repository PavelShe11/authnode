package domain

type Account struct {
	Id        string `db:"id"`
	FirstName string `db:"first_name" validate:"required,min=2,max=50"`
	LastName  string `db:"last_name" validate:"required,min=2,max=50"`
	Email     string `db:"email" validate:"required,email,min=6,max=100"`
}
