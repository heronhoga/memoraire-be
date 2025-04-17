package requests

type RequestRegister struct {
	Username  string    `json:"username" validate:"required"`
	Password  string    `json:"password" validate:"required,min=8"`
	FirstName string    `json:"first_name" validate:"required"`
	LastName  string    `json:"last_name" validate:"required"`
}

type RequestLogin struct {
	Username  string    `json:"username" validate:"required"`
	Password  string    `json:"password" validate:"required"`
}