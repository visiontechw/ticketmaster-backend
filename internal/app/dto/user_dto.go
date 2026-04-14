package dto

type CreateUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserOutput struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginOutput struct {
	Token string `json:"token"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
