package domain

type Error struct {
	Name        string
	FieldErrors []FieldError
}

type FieldError struct {
	Name    string
	Message string
}
