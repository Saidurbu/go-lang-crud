package postgres

import (
	"database/sql"
	"fmt"

	"github.com/Saidurbu/go-lang-crud/internal/config"
	"github.com/Saidurbu/go-lang-crud/internal/types"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type Postgres struct {
	DB *sql.DB
}

func New(cfg *config.Config) (*Postgres, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS students (
			id SERIAL PRIMARY KEY,
			name TEXT,
			email TEXT UNIQUE,
			password TEXT,
			age INTEGER
		)
	`)
	if err != nil {
		return nil, err
	}

	return &Postgres{DB: db}, nil
}

func (p *Postgres) CreateStudent(name string, email string, password string, age int) (int64, error) {
	if password == "" {
		return 0, fmt.Errorf("password is required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	stmt, err := p.DB.Prepare("INSERT INTO students (name, email, password, age) VALUES ($1, $2, $3, $4) RETURNING id")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var lastId int64
	err = stmt.QueryRow(name, email, string(hashedPassword), age).Scan(&lastId)
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (p *Postgres) GetStudentById(id int64) (types.Student, error) {
	stmt, err := p.DB.Prepare("SELECT id, name, email, password, age FROM students WHERE id = $1")
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
		return types.Student{}, err
	}
	return student, nil
}

func (p *Postgres) GetStudents() ([]types.Student, error) {
	stmt, err := p.DB.Prepare("SELECT id, name, email, password, age FROM students")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
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

func (p *Postgres) UpdateStudent1(id int64, name string, email string, password string, age int) error {
	stmt, err := p.DB.Prepare("UPDATE students SET name = $1, email = $2, password = $3, age = $4 WHERE id = $5")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(name, email, password, age, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("student not found with id %d: %w", id, err)
		}
		return err
	}

	return nil
}

func (p *Postgres) UpdateStudent(id int64, name string, email string, password string, age int) error {

	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		stmt, err := p.DB.Prepare("UPDATE students SET name = $1, email = $2, password = $3, age = $4 WHERE id = $5")
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(name, email, string(hashedPassword), age, id)
		if err != nil {
			return err
		}

	} else {
		stmt, err := p.DB.Prepare("UPDATE students SET name = $1, email = $2, age = $3 WHERE id = $4")
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

func (p *Postgres) DeleteStudent(id int64) error {
	stmt, err := p.DB.Prepare("DELETE FROM students WHERE id = $1")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("student not found with id %d: %w", id, err)
		}
		return err
	}

	return nil
}

func (p *Postgres) GetStudentByEmail(email string) (types.Student, error) {
	stmt, err := p.DB.Prepare("SELECT id, name, email, password, age FROM students WHERE email = $1")
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
		return types.Student{}, err
	}
	return student, nil
}
