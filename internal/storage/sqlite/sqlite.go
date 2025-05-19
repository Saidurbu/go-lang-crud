package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/Saidurbu/go-lang-crud/internal/config"
	"github.com/Saidurbu/go-lang-crud/internal/types"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type Sqlite struct {
	DB *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)

	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT,email Text,password Text, age INTEGER)`)
	if err != nil {
		return nil, err
	}

	return &Sqlite{DB: db}, nil
}

func (s *Sqlite) CreateStudent(name string, email string, password string, age int) (int64, error) {
	if password == "" {
		return 0, fmt.Errorf("password is required")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	stmt, err := s.DB.Prepare("INSERT INTO students (name, email, password, age) VALUES (?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(name, email, string(hashedPassword), age)
	if err != nil {
		return 0, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.DB.Prepare("SELECT id, name, email, password, age FROM students WHERE id = ?")
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	var student types.Student
	err = stmt.QueryRow(id).Scan(&student.ID, &student.Name, &student.Email, &student.Password, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("student not found with id %d: %w", id, err)
		}
		return types.Student{}, fmt.Errorf("query error: %w", err)
	}

	return student, nil
}

func (s *Sqlite) GetStudentByEmail(email string) (types.Student, error) {
	stmt, err := s.DB.Prepare("SELECT id, name, email, password, age FROM students WHERE email = ?")
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()
	var student types.Student
	err = stmt.QueryRow(email).Scan(&student.ID, &student.Name, &student.Email, &student.Password, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("student not found with email %s: %w", email, err)
		}
		return types.Student{}, fmt.Errorf("query error: %w", err)
	}
	return student, nil
}

func (s *Sqlite) GetStudents() ([]types.Student, error) {
	stmt, err := s.DB.Prepare("SELECT id, name, email, password, age FROM students")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var students []types.Student
	for rows.Next() {
		var student types.Student
		err = rows.Scan(&student.ID, &student.Name, &student.Email, &student.Password, &student.Age)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	return students, nil
}

func (s *Sqlite) UpdateStudent(id int64, name string, email string, password string, age int) error {

	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		stmt, err := s.DB.Prepare("UPDATE students SET name = ?, email = ?, password = ?, age = ? WHERE id = ?")
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(name, email, string(hashedPassword), age, id)
		if err != nil {
			return err
		}

	} else {
		stmt, err := s.DB.Prepare("UPDATE students SET name = ?, email = ?, age = ? WHERE id = ?")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(name, email, age, id)
		if err != nil {
			return err
		}

	}
	return nil

}

func (s *Sqlite) UpdateStudent1(id int64, name string, email string, password string, age int) error {

	var hashedPassword string

	if password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		hashedPassword = string(hashed)

		stmt, err := s.DB.Prepare("UPDATE students SET name = ?, email = ?, password = ?, age = ? WHERE id = ?")
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(name, email, hashedPassword, age, id)
		if err != nil {
			return err
		}
	} else {
		// Don't update password
		stmt, err := s.DB.Prepare("UPDATE students SET name = ?, email = ?, age = ? WHERE id = ?")
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(name, email, age, id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Sqlite) DeleteStudent(id int64) error {
	stmt, err := s.DB.Prepare("DELETE FROM students WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
