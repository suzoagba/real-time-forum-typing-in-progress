package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"real-time-forum/internal/structs"
)

func comment(r *http.Request, db *sql.DB, forPage structs.ForPage) structs.ForPage {
	var credentials struct {
		Comment string `json:"comment"`
		ID      string `json:"id"`
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
	content := credentials.Comment
	id := credentials.ID
	// Insert the reply data into the database
	_, err = db.Exec("INSERT INTO comments (postID, username, content) VALUES (?, ?, ?)", id, forPage.User.Username, content)
	if err != nil {
		return HttpError(forPage, http.StatusInternalServerError, "Failed to insert reply data into database")
	}
	// Redirect or display a success message
	return viewPost(db, forPage, id)
}
