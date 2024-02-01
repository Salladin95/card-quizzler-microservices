package main

type SignInDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"min=6,required"`
}

type SighUpDto struct {
	Name     string `json:"name"  validate:"required,min=1"`
	Password string `json:"password"  validate:"required,min=6"`
	Email    string `json:"email"  validate:"required,email"`
	Birthday string `json:"birthday"  validate:"required,min=1"`
}
