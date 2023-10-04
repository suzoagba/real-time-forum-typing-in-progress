package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"real-time-forum/internal/structs"
)

func register(w http.ResponseWriter, r *http.Request, db *sql.DB, forPage structs.ForPage) structs.ForPage {
	var credentials struct {
		Username  string `json:"username"`
		Age       string `json:"age"`
		Gender    string `json:"gender"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		return HttpError(forPage, http.StatusBadRequest, "Bad request", "Unable to decode JSON")
	}

	// Extract the user input values
	username := credentials.Username
	age := credentials.Age
	gender := credentials.Gender
	firstName := credentials.FirstName
	lastName := credentials.LastName
	email := credentials.Email
	password := credentials.Password

	invalidInput := false // Checks if the credentials that user wrote are valid
	// Check if the email or username is already taken
	if rowExists("SELECT email from users WHERE email = ?", email, db) { // if the email exists
		invalidInput = true
		forPage.Error.Error = true
		forPage.Error.Message = "Email already taken"
		return forPage
	} else if rowExists("SELECT username from users WHERE username = ?", username, db) { // if the email exists
		invalidInput = true
		forPage.Error.Error = true
		forPage.Error.Message = "Username already taken"
		return forPage
	}

	// Encrypt the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return HttpError(forPage, http.StatusInternalServerError, "Failed to encrypt password")
	}

	if !invalidInput {
		// Insert the user into the database
		err = insertUser(username, age, gender, firstName, lastName, email, string(hashedPassword), db)
		if err != nil {
			return HttpError(forPage, http.StatusInternalServerError, "Failed to register user")
		}
		return signIn(w, db, forPage, username)
	} else {
		fmt.Println("- Registration failed!")
	}
	return forPage
}

// Function to check if the email is already taken (example implementation)
func rowExists(query string, value string, db *sql.DB) bool {
	row := db.QueryRow(query, value)
	switch err := row.Scan(&value); err {
	case sql.ErrNoRows:
		return false
	case nil:
		return true
	default:
		return false
	}
}

// Function to insert the user into the database (example implementation)
func insertUser(username, age, gender, firstName, lastName, email, password string, db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO users (userID, username, email, age, gender, firstName, lastName, password) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	stmt.Exec(uuid.New().String(), username, email, age, gender, firstName, lastName, password)
	return nil
}
