package storage

type Storage interface {
	// CreateStudent creates a new student record in the storage.
	CreateStudent(name string, email string, age int) (int64, error)
}
