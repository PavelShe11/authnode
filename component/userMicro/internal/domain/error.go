package domain

type Error struct {
	Name        string
	FieldErrors []FieldError
}

type FieldError struct {
	Name    string            `json:"name"`
	Message string            `json:"message"`
	Params  map[string]string `json:"-"` // Parameters for validation errors (not serialized to JSON)
}
