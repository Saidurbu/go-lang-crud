package storage

import "github.com/Saidurbu/go-lang-crud/internal/types"

type Storage interface {
	CreateStudent(name string, email string, password string, age int) (uint, error)
	GetStudentById(id uint) (types.Student, error)
	GetStudents() ([]types.Student, error)
	UpdateStudent(id uint, name string, email string, password string, age int) error
	DeleteStudent(id uint) error
	GetStudentByEmail(email string) (types.Student, error)
}
