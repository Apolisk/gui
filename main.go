package main

import (
	"log"

	"main/database"
	"main/gui"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/pressly/goose"
)

func main() {
	db, err := database.Open()
	if err != nil {
		log.Panic(err)
	}

	if err = goose.SetDialect("sqlite3"); err != nil {
		log.Panic(errors.Wrap(err, "setting dialect"))
	}

	if err = goose.Up(db.DB, "migrations"); err != nil {
		log.Panic(err)
	}

	app := gui.New(db)
	if err = app.Build(); err != nil {
		log.Panic(err)
	}
}
