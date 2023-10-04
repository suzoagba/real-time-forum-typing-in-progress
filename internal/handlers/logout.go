package handlers

import (
	"database/sql"
	"net/http"
	"real-time-forum/internal/structs"
)

func logout(w http.ResponseWriter, db *sql.DB, pageData structs.ForPage) structs.ForPage {
	deleteCookie(w, "forum-session2")
	err := deleteUserSession(db, pageData.User.Username)
	if err != nil {
		return HttpError(pageData, http.StatusInternalServerError, "Failed to delete user session")
	}
	pageData = structs.ForPage{}
	pageData.User.LoggedIn = false
	return pageData
}

// Delete session cookie
func deleteCookie(w http.ResponseWriter, cookieName string) {
	http.SetCookie(w, &http.Cookie{
		Name:   cookieName,
		Value:  "",
		MaxAge: -1, // Set MaxAge to a negative value to delete the cookie
	})
}

// Delete session from DB
func deleteUserSession(db *sql.DB, username string) error {
	deleteSQL := "DELETE FROM authenticated_users WHERE username = ?"

	_, err := db.Exec(deleteSQL, username)
	return err
}
