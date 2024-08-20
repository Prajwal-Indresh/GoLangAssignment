package student

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type Student struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Age       int       `db:"age" json:"age"`
	CreatedBy string    `db:"created_by" json:"created_by"`
	CreatedOn time.Time `db:"created_on" json:"created_on"`
	UpdatedBy string    `db:"updated_by" json:"updated_by"`
	UpdatedOn time.Time `db:"updated_on" json:"updated_on"`
}

type Service struct {
	DB *sqlx.DB
}

func NewService(db *sqlx.DB) *Service {
	return &Service{DB: db}
}

// CreateStudent
func (s *Service) CreateStudent(student *Student) error {
	student.CreatedOn = time.Now()
	student.UpdatedOn = time.Now()
	_, err := s.DB.NamedExec(`INSERT INTO students (name, age, created_by, created_on, updated_by, updated_on) 
	VALUES (:name, :age, :created_by, :created_on, :updated_by, :updated_on)`, student)
	if err != nil {
		log.Printf("Error inserting student: %v", err)
	}
	return err
}

// GetStudent By ID
func (s *Service) GetStudentByID(id int64) (*Student, error) {
	var student Student
	err := s.DB.Get(&student, "SELECT * FROM students WHERE id = ?", id)
	return &student, err
}

// UpdateStudent
func (s *Service) UpdateStudent(student *Student) error {
	student.UpdatedOn = time.Now()
	_, err := s.DB.NamedExec(`UPDATE students SET name = :name, age = :age, updated_by = :updated_by, updated_on = :updated_on WHERE id = :id`, student)
	return err
}

// DeleteStudent By ID
func (s *Service) DeleteStudent(id int64) error {
	_, err := s.DB.Exec("DELETE FROM students WHERE id = ?", id)
	return err
}

// GetAllStudents
func (s *Service) GetAllStudents() ([]Student, error) {
	var students []Student
	err := s.DB.Select(&students, "SELECT * FROM students")
	if err != nil {
		log.Printf("Error retrieving students: %v", err)
	}
	return students, err
}
