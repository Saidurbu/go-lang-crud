package storage

import "github.com/Saidurbu/go-lang-crud/internal/types"

type Storage interface {
	CreateStudent(name string, email string, password string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetStudents() ([]types.Student, error)
	UpdateStudent(id int64, name string, email string, password string, age int) error
	DeleteStudent(id int64) error
	GetStudentByEmail(email string) (types.Student, error)
}
