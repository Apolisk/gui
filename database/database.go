package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	*sqlx.DB
}

func Open(url string) (*DB, error) {
	db, err := sqlx.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &DB{DB: db}, nil
}

func (db *DB) HasUser(username string) (has bool) {
	const query = `SELECT EXISTS (
		SELECT 1 FROM information WHERE username = $1);`
	_ = db.Get(&has, query, username)
	return has
}

func (db *DB) CreateUser(username, password string) error {
	const query = `INSERT INTO information values($1, $2)`
	_, err := db.Exec(query, username, password)
	return err

}

func (db *DB) ChangePassword(username, password string) error {
	const query = `UPDATE information SET user_password = $2 WHERE username = $1`
	_, err := db.Exec(query, username, password)
	return err
}

func (db *DB) LogIn(username, password string) (has bool) {
	const query = `SELECT EXISTS (
		SELECT 1 FROM information WHERE username = $1 AND user_password = $2);`
	_ = db.Get(&has, query, username, password)
	return has
}
