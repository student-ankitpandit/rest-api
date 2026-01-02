package storage

import "github.com/student-ankitpandit/rest-api/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentsById() (types.Student, error)
	GetStudentsList ([]types.Student, error)
}