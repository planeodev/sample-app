// order-service/main.go
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

type Order struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

func connectDB() {
	var err error
	// Get database credentials from environment variables
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("PSQL_DB_USER"), os.Getenv("PSQL_DB_PASSWORD"),
		os.Getenv("PSQL_DB_HOST"), os.Getenv("PSQL_DB_PORT"), os.Getenv("PSQL_DB_NAME"))
	log.Default().Printf("db conn string is %s\n", connStr)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

func GetOrders(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, user_id, product_id, quantity FROM orders")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.ID, &order.UserID, &order.ProductID, &order.Quantity); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		orders = append(orders, order)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func main() {
	log.Default().Println("Starting order service...")
	connectDB()
	log.Default().Println("Database connected..")

	defer db.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	r := mux.NewRouter()

	r.Handle("/", http.FileServer(http.Dir("./static")))
	r.HandleFunc("/orders", GetOrders).Methods("GET")
	log.Default().Printf("Serving API on port %s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}
