package database

import (
	"log"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func Open() (*DB, error) {
	db, err := sql.Open("sqlite3", "database/info.db")
	if err != nil {
		return nil, err
	}
	return &DB{DB: db}, nil
}

func (db *DB) HasUser(username string) bool {
	const query = `SELECT COUNT(*) FROM users WHERE name = $1;`
	var count int
	if err := db.QueryRow(query, username).Scan(&count); err != nil {
		log.Print(err)
		return false
	}
	return count > 0
}

func (db *DB) CreateUser(username, password string) error {
	const query = `INSERT INTO users values($1, $2, 0);`
	_, err := db.Exec(query, username, password)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) ChangePassword(password, username string) error {
	const query = `UPDATE users SET password = $1 WHERE name = $2;`
	_, err := db.Exec(query, password, username)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) LogIn(username, password string) bool {
	const query = `SELECT COUNT(*) FROM users WHERE name = $1 AND password = $2;`
	var count int
	if err := db.QueryRow(query, username, password).Scan(&count); err != nil {
		log.Print(err)
		return false
	}
	return count > 0
}

func (db *DB) Level(username string) (int, error) {
	const query = `SELECT access_level from users where name = $1;`
	row := db.QueryRow(query, username)
	var accessLevel int
	err := row.Scan(&accessLevel)
	if err != nil {
		return -1, err
	}
	return accessLevel, nil
}
