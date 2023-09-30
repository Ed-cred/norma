package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "rootuser"
	dbname   = "norma"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	// db, err := sql.Open("postgres", psqlInfo)
	// if err != nil {
	// 	panic(err)
	// }
	// err = resetDB(db, dbname)
	// if err != nil {
	// 	panic(err)
	// }
	// db.Close()
	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// err = db.Ping()
	// if err != nil {
	// 	log.Println("couldn't ping the database:", err)
	// }
	err = createPhoneNumbersTable(db)
	if err != nil {
		log.Println("couldn't create phone numbers table:", err)
	}
	_, err = insertPhone(db, "1234567890")
	if err != nil {
		log.Println("unable to insert phone numberL", err)
	}
	_, err = insertPhone(db, "123 456 7891")
	if err != nil {
		log.Println("unable to insert phone numberL", err)
	}
	_, err = insertPhone(db, "(123) 456 7892")
	if err != nil {
		log.Println("unable to insert phone numberL", err)
	}
	_, err = insertPhone(db, "(123) 456-7893")
	if err != nil {
		log.Println("unable to insert phone numberL", err)
	}
	_, err = insertPhone(db, "123-456-7894")
	if err != nil {
		log.Println("unable to insert phone numberL", err)
	}
	_, err = insertPhone(db, "123-456-7890")
	if err != nil {
		log.Println("unable to insert phone numberL", err)
	}
	_, err = insertPhone(db, "1234567892")
	if err != nil {
		log.Println("unable to insert phone numberL", err)
	}
	_, err = insertPhone(db, "(123)456-7892")
	if err != nil {
		log.Println("unable to insert phone numberL", err)
	}

	// id, err := insertPhone(db, "1234567890")
	// if err != nil {
	// 	log.Println("unable to insert phone numberL", err)
	// }
	// log.Println("created phone number record:", id)
}

func insertPhone(db *sql.DB, phone string) (int, error) {
	stmt := `INSERT INTO phone_numbers(value) VALUES($1) RETURNING id`
	var id int
	err := db.QueryRow(stmt, phone).Scan(&id)
	if err != nil {
		log.Println("unable to insert phone number into db:", err)
		return -1, err
	}
	return id, nil
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

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
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

// iter version
func normalize(phone string) string {
	var buf bytes.Buffer
	for _, ch := range phone {
		if ch >= '0' && ch <= '9' {
			buf.WriteRune(ch)
		}
	}
	return buf.String()
}

// regex version
// func normalize(phone string) string {
// 	re := regexp.MustCompile("[^0-9]")
// 	return re.ReplaceAllString(phone, "")
// }
