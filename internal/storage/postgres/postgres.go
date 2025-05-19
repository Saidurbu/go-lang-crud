package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/Saidurbu/go-lang-crud/internal/config"
	"github.com/Saidurbu/go-lang-crud/internal/types"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	DB *gorm.DB
}

func New(cfg *config.Config) (*Postgres, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&types.Student{})
	if err != nil {
		return nil, err
	}
	log.Println("GORM connected to DB")
	return &Postgres{DB: db}, nil
}

func (p *Postgres) CreateStudent(name string, email string, password string, age int) (uint, error) {
	if password == "" {
		return 0, fmt.Errorf("password is required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	student := types.Student{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Age:      age,
	}

	if err := p.DB.Create(&student).Error; err != nil {
		return 0, err
	}

	return student.ID, nil
}

func (p *Postgres) GetStudentById(id uint) (types.Student, error) {
	var student types.Student
	if err := p.DB.First(&student, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return types.Student{}, fmt.Errorf("student not found with id %d: %w", id, err)
		}
		return types.Student{}, fmt.Errorf("query error: %w", err)
	}
	return student, nil
}

func (p *Postgres) GetStudents() ([]types.Student, error) {
	var students []types.Student
	if err := p.DB.Find(&students).Error; err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	return students, nil
}

func (p *Postgres) GetStudentByEmail(email string) (types.Student, error) {
	stmt, err := p.DB.Raw("SELECT id, name, email, password, age FROM students WHERE email = $1", email).Rows()
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	var student types.Student
	if stmt.Next() {
		if err := stmt.Scan(&student.ID, &student.Name, &student.Email, &student.Password, &student.Age); err != nil {
			return types.Student{}, err
		}
	} else {
		return types.Student{}, fmt.Errorf("student not found with email %s: %w", email, sql.ErrNoRows)
	}

	return student, nil
}

func (p *Postgres) UpdateStudent1(id uint, name string, email string, password string, age int) error {
	stmt, err := p.DB.Raw("UPDATE students SET name = $1, email = $2, password = $3, age = $4 WHERE id = $5", name, email, password, age, id).Rows()
	if err != nil {
		return err
	}
	defer stmt.Close()
	// Check if the student exists
	var student types.Student
	err = stmt.Scan(&student.ID, &student.Name, &student.Email, &student.Password, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("student not found with id %d: %w", id, err)
		}
		return fmt.Errorf("query error: %w", err)
	}
	return nil
}

func (p *Postgres) UpdateStudent(id uint, name, email, password string, age int) error {
	var student types.Student

	if err := p.DB.First(&student, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("student not found with id %d", id)
		}
		return fmt.Errorf("failed to find student: %w", err)
	}

	student.Name = name
	student.Email = email
	student.Age = age

	if password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		student.Password = string(hashed)
	}

	if err := p.DB.Save(&student).Error; err != nil {
		return fmt.Errorf("failed to update student: %w", err)
	}

	return nil
}

func (p *Postgres) DeleteStudent(id uint) error {
	var student types.Student
	result := p.DB.First(&student, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("student not found with id %d", id)
		}
		return result.Error
	}

	if err := p.DB.Delete(&student).Error; err != nil {
		return fmt.Errorf("failed to delete student: %w", err)
	}

	return nil
}
