package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type DataItem struct {
	Label string `json:"label"`
	Value int    `json:"value"`
}

var db *sql.DB

func main() {

	// Initialize SQLite database
	var err error
	db, err = sql.Open("sqlite3", "/Users/bfrankl/Desktop/project/MyMicro/check/db/mydatabase.db")

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/api/data/table", handleTableData)
	http.HandleFunc("/api/data/pie", handlePieChartData)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/signup", handleSignup)
	http.Handle("/", http.FileServer(http.Dir(".")))

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleTableData(w http.ResponseWriter, r *http.Request) {
	// Sample table data
	data := []DataItem{
		{Label: "Item 1", Value: 10},
		{Label: "Item 2", Value: 20},
		{Label: "Item 3", Value: 30},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func handlePieChartData(w http.ResponseWriter, r *http.Request) {
	// Sample pie chart data
	data := []DataItem{
		{Label: "Category 1", Value: 50},
		{Label: "Category 2", Value: 30},
		{Label: "Category 3", Value: 20},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func handleSignup(w http.ResponseWriter, r *http.Request) {
	// Extract username and password from the request body
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Println(username)
	fmt.Println(password)
	// Insert user data into the database
	db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, password)
	_, err := db.Exec(`INSERT INTO users (username, password) VALUES (?, ?)`, "admin1", "admin1")
	fmt.Println(err)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User signup successful"))
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	// Extract username and password from the request body
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Query the database to check if the user exists
	var storedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Compare the stored password with the provided password
	if password != storedPassword {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Send a success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User login successful"))
}
