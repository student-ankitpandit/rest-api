package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/student-ankitpandit/rest-api/internal/config"
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
	stmt, err := s.Db.Prepare("INSERTS INTO STUDENTS (name, email, password) VALUES (?, ?, ?)")
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