// product-service/main.go
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

type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
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

func GetProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, price FROM products")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func main() {
	log.Default().Println("Starting product service...")
	connectDB()
	log.Default().Println("Database connected..")
	defer db.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	r := mux.NewRouter()

	r.Handle("/", http.FileServer(http.Dir("./static")))
	r.HandleFunc("/products", GetProducts).Methods("GET")
	log.Default().Printf("Serving API on port %s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}
