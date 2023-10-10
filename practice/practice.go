package main

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

func main() {
	db, err := sql.Open("postgres", "sslmode=disable dbname=test user=test password=test host=127.0.0.1 port=9393")
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Print("database connected successfully!!!")
	fmt.Println("acra practice")

	http.HandleFunc("/create-table", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("handler called")
		_, err := db.Exec("CREATE TABLE IF NOT EXISTS users(id SERIAL PRIMARY KEY ,name TEXT,email TEXT,password TEXT);")
		if err != nil {
			log.Fatal(err)
		}

	})
	http.HandleFunc("/add-data", func(w http.ResponseWriter, r *http.Request) {
		_, err := db.Exec(`INSERT INTO  users(id,name,email,password) VALUES(1,'test','test@gmail.com','test@123');`)
		if err != nil {
			log.Fatal(err)
		}

	})
	http.HandleFunc("/get-data", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * from users")
		if err != nil {
			log.Fatal(err)
		}
		type Row struct {
			ID       int
			Name     string
			Email    string
			Password string
		}
		for rows.Next() {
			var row Row
			err := rows.Scan(&row.ID, &row.Name, &row.Email, &row.Password)

			if err != nil {
				log.Fatal(err)
			}

			// // Remove escape characters and decode hexadecimal strings
			// cleanedEmail := strings.ReplaceAll(row.Email, "\\x", "")
			// cleanedPassword := strings.ReplaceAll(row.Password, "\\x", "")

			// decodedEmail, err := hex.DecodeString(cleanedEmail)
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// decodedPassword, err := hex.DecodeString(cleanedPassword)
			// if err != nil {
			// 	log.Fatal(err)
			// }

			fmt.Println(row.ID, " ", row.Name, " ", row.Email," ",row.Password)
		}
	})

	http.HandleFunc("/get-one",func(w http.ResponseWriter,r *http.Request){
		rows,err:=db.Query(`select * from users where name='test'`)
		if err!=nil{
			log.Fatal(err)
		}
		
		for rows.Next(){
			var row User

			err:=rows.Scan(&row.ID,&row.Name,&row.Email,&row.Password)
			if err!=nil{
				log.Fatal(err)
			}
			cleanedEmail := strings.ReplaceAll(row.Email, "\\x", "")
			cleanedPassword := strings.ReplaceAll(row.Password, "\\x", "")

			decodedEmail, err := hex.DecodeString(cleanedEmail)
			if err != nil {
				log.Fatal(err)
			}
			decodedPassword, err := hex.DecodeString(cleanedPassword)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(row.ID, " ", row.Name, " ", string(decodedEmail), " ", string(decodedPassword))



		}
	})


	log.Fatal(http.ListenAndServe(":3000", nil))
}
