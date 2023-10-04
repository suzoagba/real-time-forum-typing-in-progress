package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"real-time-forum/internal/structs"
)

func createPost(r *http.Request, db *sql.DB, forPage structs.ForPage) structs.ForPage {
	fmt.Println("[ CreatePost ] function")
	var credentials struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		return HttpError(forPage, http.StatusBadRequest, "Bad request", err.Error())
	}

	err := r.ParseForm()
	if err != nil {
		forPage.Error.Error = true
		forPage.Error.Message = "Failed to parse form data"
		return forPage
	}

	// Extract the user input values
	title := credentials.Title
	description := credentials.Description
	selectedTags := credentials.Tags

	fmt.Println(title, description, selectedTags)
	// Check if data includes empty fields
	if title == "" || description == "" || len(selectedTags) == 0 {
		return HttpError(forPage, http.StatusBadRequest, "Forbidden empty fields.")
	}
	// Insert the post data into the database
	result, err := db.Exec("INSERT INTO posts (username, title, description) VALUES (?, ?, ?)", forPage.User.Username, title, description)
	if err != nil {
		return HttpError(forPage, http.StatusInternalServerError, "Failed to insert post data into database.")
	}
	// Get the postID of the newly created post
	postID, err := result.LastInsertId()
	if err != nil {
		return HttpError(forPage, http.StatusInternalServerError, "Failed to get post ID.")
	}
	// Process the selected tags and insert into post_tags table
	for _, tagID := range selectedTags {
		_, err = db.Exec("INSERT INTO post_tags (postID, tagID) VALUES (?, ?)", postID, tagID)
		if err != nil {
			return HttpError(forPage, http.StatusInternalServerError, "Failed to insert tag into post_tags table")
		}
	}
	return forPage
}
