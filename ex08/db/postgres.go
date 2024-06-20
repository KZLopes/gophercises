package db

import (
	"database/sql"
	"fmt"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "carlos"
	password = "123456"
	dbname   = "gophercises_phone"
)

func InitDd() (*sql.DB, error) {
	// psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)

	// 	db, err := sql.Open("postgres", psqlInfo)
	// 	must(err)
	//
	// 	err = resetDB(db, dbname)
	// 	must(err)
	//
	// 	db.Close()
	//
	// 	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name) //
	if err != nil {
		return err
	}

	return nil
}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return err
	}

	return createDB(db, name)
}

func createPhoneNumberTable(db *sql.DB) error {
	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS phone_number (
			id SERIAL,
			value VARCHAR(255),
			normalized VARCHAR(255)
		)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}
