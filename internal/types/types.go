package types

type Student struct {
	ID       uint `gorm:"primaryKey"`
	Name     string
	Email    string `gorm:"unique"`
	Password string
	Age      int
}

type StudentResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}
