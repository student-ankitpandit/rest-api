package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/student-ankitpandit/rest-api/internal/config"
	"github.com/student-ankitpandit/rest-api/internal/types"
)

type Sqlite struct {
	Db *sql.DB
}

//making an instace of Sqlite struct
func New(cfg *config.Config) (*Sqlite, error) {
	//opening a connection to sqlite database
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if(err != nil) {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT, 
	age INTEGER
	)`)

	if(err != nil) {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil //mutiple values returning here (db connection and nil)
}

func (s *Sqlite)CreateStudent(name string, email string, age int) (int64, error)  {
	stmt, err := s.Db.Prepare("INSERT INTO STUDENTS (name, email, age) VALUES (?, ?, ?)")
	if(err != nil) {
		return 0, err
	}

	defer stmt.Close()

	res, err := stmt.Exec(name, email, age)
	if(err != nil) {
		return 0, err
	}

	lastId, err := res.LastInsertId()
	if(err != nil) {
		return 0, err
	}

	return lastId, nil
}

func (s * Sqlite) GetStudentsById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students WHERE id = ? LIMIT 1")
	if(err != nil) {
		return types.Student{}, err
	}

	defer stmt.Close()

	//serializing the data into a struct

	var student types.Student

	err = stmt.QueryRow(id).Scan(student.Id, student.Name, student.Email, student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no user found with id %s", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf("query error %w", err)
	}
}