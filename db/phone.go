package db

import (
	"database/sql"
)

type DB struct {
	db *sql.DB
}

type Phone struct {
	ID    int
	Value string
}

func Open(driverName, dataSource string) (*DB, error) {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (d *DB) Close() error {
	return d.db.Close()
}

func (d *DB) Seed() error {
	data := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}
	for _, number := range data {
		if _, err := insertPhone(d.db, number); err != nil {
			return err
		}
	}
	return nil
}

func Reset(driverName, dataSource, dbName string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}
	err = resetDB(db, dbName)
	if err != nil {
		return err
	}
	return db.Close()
}

// Migrate handles creating the phone table and adding data to it
func Migrate(driverName, dataSource string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}
	err = createPhoneNumbersTable(db)
	if err != nil {
		return err
	}
	return db.Close()
}

func (d *DB) AllPhones() ([]Phone, error) {
	rows, err := d.db.Query("SELECT id, value FROM phone_numbers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ret []Phone
	for rows.Next() {
		var p Phone
		if err := rows.Scan(&p.ID, &p.Value); err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}

func (d *DB) DeletePhone(id int) error {
	stmt := `DELETE FROM phone_numbers WHERE id=$1`
	_, err := d.db.Exec(stmt, id)
	return err
}

func (d *DB) UpdatePhone(p Phone) error {
	stmt := `UPDATE phone_numbers SET value = $2 WHERE id=$1`
	_, err := d.db.Exec(stmt, p.ID, p.Value)
	return err
}

func (d *DB) FindPhone(number string) (*Phone, error) {
	var p Phone
	row := d.db.QueryRow("SELECT id, value FROM phone_numbers WHERE value=$1", number)
	err := row.Scan(&p.ID, &p.Value)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &p, nil
}

func createPhoneNumbersTable(db *sql.DB) error {
	stmt := `
		CREATE TABLE IF NOT EXISTS phone_numbers (
			id SERIAL,
			value  VARCHAR(255)
		)`
	_, err := db.Exec(stmt)
	return err
}

func insertPhone(db *sql.DB, phone string) (int, error) {
	stmt := `INSERT INTO phone_numbers(value) VALUES($1) RETURNING id`
	var id int
	err := db.QueryRow(stmt, phone).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name + "WITH (FORCE)")
	if err != nil {
		return err
	}
	return createDB(db, name)
}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	if err != nil {
		return err
	}
	return nil
}
