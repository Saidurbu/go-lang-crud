package types

type Student struct {
	ID       int64  `json:"id"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" min:"8"`
	Age      int    `json:"age" validate:"required,min=0"`
}

type StudentResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}
