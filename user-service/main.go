// user-service/main.go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func connectDB() {
	var err error
	// Get database credentials from environment variables
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("PSQL_DB_USER"), os.Getenv("PSQL_DB_PASSWORD"),
		os.Getenv("PSQL_DB_HOST"), os.Getenv("PSQL_DB_PORT"), os.Getenv("PSQL_DB_NAME"))

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func main() {
	log.Default().Println("Starting user service...")
	connectDB()
	log.Default().Println("Database connected..")

	defer db.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r := mux.NewRouter()

	r.Handle("/", http.FileServer(http.Dir("./static")))
	r.HandleFunc("/users", GetUsers).Methods("GET")
	log.Default().Printf("Serving API on port %s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}
