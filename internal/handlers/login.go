package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"real-time-forum/internal/structs"
	"time"
)

func login(w http.ResponseWriter, r *http.Request, db *sql.DB, forPage structs.ForPage) structs.ForPage {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		return HttpError(forPage, http.StatusBadRequest, "Bad request")
	}

	err := r.ParseForm()
	if err != nil {
		forPage.Error.Error = true
		forPage.Error.Message = "Failed to parse form data"
		return forPage
	}

	// Extract the user input values
	username := credentials.Username
	password := credentials.Password
	if passwordCorrect("SELECT password FROM users WHERE username = ?", username, password, db) {
		return signIn(w, db, forPage, username)
	} else if passwordCorrect("SELECT password FROM users WHERE email = ?", username, password, db) {
		var realUsername string
		db.QueryRow("SELECT username FROM users WHERE email = ?", username).Scan(&realUsername)
		return signIn(w, db, forPage, realUsername)
	} else {
		forPage.Error.Error = true
		forPage.Error.Message = "Check your password"
		return forPage
	}
}
func updateSessionID(uuid, sessionID string, db *sql.DB) error {
	// Prepare the SQL statement
	stmt, err := db.Prepare("UPDATE authenticated_users SET sessionID = ? WHERE username = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	// Execute the SQL statement
	_, err = stmt.Exec(sessionID, uuid)
	if err != nil {
		return err
	}
	return nil
}
func passwordCorrect(query string, value string, password string, db *sql.DB) bool {
	row := db.QueryRow(query, value)
	var storedPassword string
	// Scan the retrieved password into the variable
	err := row.Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			// User not found in the database
			return false
		}
		return false
	}
	// Compare the stored password with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		// Password does not match
		return false
	}
	return true
}

func signIn(w http.ResponseWriter, db *sql.DB, forPage structs.ForPage, username string) structs.ForPage {
	sessionID := uuid.New().String()
	expiration := time.Now().Add(24 * time.Hour) // Set the expiration time for the cookie
	createCookie(w, "forum-session2", sessionID, expiration)
	err := updateSessionID(username, sessionID, db)
	if err != nil {
		forPage.Error.Error = true
		forPage.Error.Message = "Unable to update session ID: " + err.Error() + "."
		return forPage
	}
	err = addActiveSession(db, sessionID, username)
	if err != nil {
		forPage.Error.Error = true
		forPage.Error.Message = "Unable to add active session"
		return forPage
	}
	forPage.User.LoggedIn = true
	forPage.User.Username = username
	forPage.User.Session = sessionID

	return forPage
}

// Create a new session cookie
func createCookie(w http.ResponseWriter, name, value string, expiration time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:   name,
		Value:  "",
		MaxAge: -1, // Set MaxAge to a negative value to delete the cookie
	})
	cookie := &http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expiration,
	}
	http.SetCookie(w, cookie)
}

// Add new active session to the database
func addActiveSession(db *sql.DB, sessionID, username string) error {
	insertSQL := `
        INSERT OR REPLACE INTO authenticated_users (sessionID, username)
        VALUES (?, ?);
    `
	_, err := db.Exec(insertSQL, sessionID, username)
	if err != nil {
		return err
	}
	return nil
}
