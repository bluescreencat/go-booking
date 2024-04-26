package dto

type LoginDTO struct {
	Username string `json:"username" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=16,password-format"`
}

// ,password-format
