package dto

type RegisterDTO struct {
	Username string `json:"username" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=16,password-format"`
	Name     string `json:"name" validate:"required,max=255"`
	Surname  string `json:"surname" validate:"required,max=255"`
}
