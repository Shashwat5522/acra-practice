package main

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

type User struct {
	ID       int
	Name     string 	`json:"name"`
	Email    string		`json:"email"`
	Password string		`json:"password"`
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
		var user User

    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        log.Fatal(err)
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Use a prepared statement to insert data into the database
    InsertQuery := "INSERT INTO users(name, email, password) VALUES($1, $2, $3);"

    // Execute the prepared statement with user data as parameters
    _, ierr := db.Exec(InsertQuery, user.Name, user.Email, user.Password)
    if ierr != nil {
        log.Fatal(ierr)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    // Send a success response to the client
    w.WriteHeader(http.StatusCreated)
    fmt.Fprintln(w, "Data successfully inserted into the database")

	})
	http.HandleFunc("/get-data", func(w http.ResponseWriter, r *http.Request) {
		var users []User
		rows, err := db.Query("SELECT * from users")
		if err != nil {
			log.Fatal(err)
		}
		
		for rows.Next() {
			var user User
			err := rows.Scan(&user.ID,&user.Name, &user.Email, &user.Password)


			if err != nil {
				log.Fatal(err)
			}
			users=append(users, user)
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
			
				
		}
		fmt.Print(users)
				w.WriteHeader(200)
				fmt.Fprint(w,users)
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
