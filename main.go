package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"real-time-forum/internal/database"
	"real-time-forum/internal/handlers"
	"real-time-forum/internal/structs"
)

func main() {
	// Create new empty database
	database.StartDB(false)
	// Open a connection to the database
	db, err := sql.Open("sqlite3", "./internal/database/forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// Create the tables if they don't exist
	err = database.CreateTables(db)
	if err != nil {
		log.Fatal(err)
	}
	structs.GetTags(db)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})
	http.HandleFunc("/page", handlers.PageHandler(db))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handlers.WsHandler(w, r, db) // Pass the response writer, request, and the database connection.
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Server listening on http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
