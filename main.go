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
	// err = createPhoneNumbersTable(db)
	// if err != nil {
	// 	log.Println("couldn't create phone numbers table:", err)
	// }
	// _, err = insertPhone(db, "1234567890")
	// if err != nil {
	// 	log.Println("unable to insert phone numberL", err)
	// }
	// _, err = insertPhone(db, "123 456 7891")
	// if err != nil {
	// 	log.Println("unable to insert phone numberL", err)
	// }
	// _, err = insertPhone(db, "(123) 456 7892")
	// if err != nil {
	// 	log.Println("unable to insert phone numberL", err)
	// }
	// id, err := insertPhone(db, "(123) 456-7893")
	// if err != nil {
	// 	log.Println("unable to insert phone numberL", err)
	// }
	// phone, err := getPhone(db, id)
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println(phone)
	// _, err = insertPhone(db, "123-456-7894")
	// if err != nil {
	// 	log.Println("unable to insert phone numberL", err)
	// }
	// _, err = insertPhone(db, "123-456-7890")
	// if err != nil {
	// 	log.Println("unable to insert phone numberL", err)
	// }
	// _, err = insertPhone(db, "1234567892")
	// if err != nil {
	// 	log.Println("unable to insert phone numberL", err)
	// }
	// _, err = insertPhone(db, "(123)456-7892")
	// if err != nil {
	// 	log.Println("unable to insert phone numberL", err)
	// }

	// // id, err := insertPhone(db, "1234567890")
	// // if err != nil {
	// // 	log.Println("unable to insert phone numberL", err)
	// // }
	// // log.Println("created phone number record:", id)
	phones, err := allPhones(db)
	if err != nil {
		log.Println(err)
	}
	for _, p := range phones {
		fmt.Printf("working on... %+v\n", p)
		number := normalize(p.value)
		if number != p.value {
			fmt.Println("Updating phone number...", number)
			existing, err := findPhone(db, number)
			if err != nil {
				log.Fatal(err)
			}
			if existing != nil {
				// delete numberr
				err := deletePhone(db, p.id)
				if err != nil {
					log.Fatal("could not delete phone number:", err)
				}
			} else {
				///update
				p.value = number
				err := updatePhone(db, p)
				if err != nil {
					log.Fatal("could not update phone number:", err)
				}
			}
		} else {
			fmt.Println("No changes required")
		}
	}
}

func deletePhone(db *sql.DB, id int) error {
	stmt := `DELETE FROM phone_numbers WHERE id=$1`
	_, err := db.Exec(stmt, id)
	return err
}

func getPhone(db *sql.DB, id int) (string, error) {
	var phone string
	err := db.QueryRow("SELECT value FROM phone_numbers WHERE id=$1", id).Scan(&phone)
	if err != nil {
		log.Println("unable to get phone number:", err)
		return "", err
	}
	return phone, nil
}

func updatePhone(db *sql.DB, p phone) error {
	stmt := `UPDATE phone_numbers SET value = $2 WHERE id=$1`
	_, err := db.Exec(stmt, p.id, p.value)
	return err
}

func findPhone(db *sql.DB, number string) (*phone, error) {
	var p phone
	row := db.QueryRow("SELECT id, value FROM phone_numbers WHERE value=$1", number)
	err := row.Scan(&p.id, &p.value)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &p, nil
}

type phone struct {
	id    int
	value string
}

func allPhones(db *sql.DB) ([]phone, error) {
	rows, err := db.Query("SELECT id, value FROM phone_numbers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ret []phone
	for rows.Next() {
		var p phone
		if err := rows.Scan(&p.id, &p.value); err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return ret, nil
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
