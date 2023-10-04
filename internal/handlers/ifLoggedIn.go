package handlers

import (
	"database/sql"
	"net/http"
	"real-time-forum/internal/structs"
)

type userInfo struct {
	User structs.User
}

func IsLoggedIn(r *http.Request, db *sql.DB) userInfo {
	info := userInfo{}
	cookie, err := r.Cookie("forum-session2")
	if err != nil {
		info.User.LoggedIn = false
		// Session ID cookie not found
		return info
	}
	sessionID := cookie.Value

	if sessionID == "" {
		info.User.LoggedIn = false
		return info
	}
	// Check if the session ID exists in the database
	row := db.QueryRow("SELECT username FROM authenticated_users WHERE sessionID = ?;", sessionID)
	var username string
	err = row.Scan(&username)
	if err == sql.ErrNoRows {
		// Session ID does not exist in the database
		info.User.LoggedIn = false
		return info
	} else if err != nil {
		// Error occurred while querying the database
		info.User.LoggedIn = false
		return info
	}
	// Check if the session ID exists in the database
	row = db.QueryRow("SELECT userID FROM users WHERE username = ?;", username)
	var uuid string
	err = row.Scan(&uuid)
	if err == sql.ErrNoRows {
		return info
	} else if err != nil {
		return info
	}
	info.User.ID = uuid
	info.User.Username = username
	info.User.LoggedIn = true
	info.User.Session = sessionID
	return info
}
