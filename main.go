package main

import (
	"bytes"
	"fmt"
	"log"

	phonedb "github.com/Ed-cred/norma/db"
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
	// err := phonedb.Reset("postgres", psqlInfo, dbname)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	// err := phonedb.Migrate("postgres", psqlInfo)
	// if err != nil {
	// 	panic(err)
	// }
	db, err := phonedb.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err := db.Seed(); err != nil {
		log.Fatal(err)
	}
	phones, err := db.AllPhones()
	if err != nil {
		log.Println(err)
		return
	}
	for _, p := range phones {
		fmt.Printf("working on... %+v\n", p)
		number := normalize(p.Value)
		if number != p.Value {
			fmt.Println("Updating phone number...", number)
			existing, err := db.FindPhone(number)
			if err != nil {
				log.Fatal(err)
			}
			if existing != nil {
				// delete numberr
				err := db.DeletePhone(p.ID)
				if err != nil {
					log.Fatal("could not delete phone number:", err)
				}
			} else {
				///update
				p.Value = number
				err := db.UpdatePhone(p)
				if err != nil {
					log.Fatal("could not update phone number:", err)
				}
			}
		} else {
			fmt.Println("No changes required")
		}
	}
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
