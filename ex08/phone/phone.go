package phone

import (
	"database/sql"
	"fmt"
)

type Phone struct {
	Id         int
	Value      string
	Normalized *string
}

func InsertPhone(db *sql.DB, phone string) (int, error) {
	stmt, err := db.Prepare(`INSERT INTO phone_number (value) VALUES ($1) RETURNING id`)
	if err != nil {
		return -1, err
	}

	var id int
	err = stmt.QueryRow(phone).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func FindPhone(db *sql.DB, id int) (*Phone, error) {
	var p Phone
	err := db.QueryRow(`SELECT * FROM phone_number WHERE id=$1`, id).Scan(&p.Id, &p.Value, &p.Normalized)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func ListPhones(db *sql.DB) ([]Phone, error) {
	var phones []Phone
	rows, err := db.Query(`SELECT * FROM phone_number`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p Phone
		err := rows.Scan(&p.Id, &p.Value, &p.Normalized)
		if err != nil {
			fmt.Println("error:", err)
			continue
		}
		phones = append(phones, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return phones, nil
}

func UpdateNormalizedNumber(db *sql.DB, p Phone) error {
	_, err := db.Exec(`UPDATE phone_number SET normalized=$1 WHERE id=$2`, p.Normalized, p.Id)
	if err != nil {
		return err
	}
	return nil
}
