package handlers

import (
	"database/sql"
	"net/http"
	"real-time-forum/internal/structs"
)

func viewPost(db *sql.DB, forPage structs.ForPage, nr string) structs.ForPage {
	// Extract the postID from the URL parameters
	postID := nr
	// Query the database to get the post information
	postQuery := `
			SELECT postID, title, description, creationDate, username
			FROM posts
			WHERE postID = ?
		`
	postRow := db.QueryRow(postQuery, postID)
	var post structs.Post
	err := postRow.Scan(&post.ID, &post.Title, &post.Description, &post.CreationDate, &post.Username)
	if err != nil {
		return HttpError(forPage, http.StatusInternalServerError, "Failed to retrieve post")
	}
	// Query the database to get the comments for the post
	commentQuery := `
			SELECT commentID, content, creationDate, username
			FROM comments
			WHERE postID = ?
			ORDER BY creationDate ASC
		`
	commentRows, err := db.Query(commentQuery, postID)
	if err != nil {
		return HttpError(forPage, http.StatusInternalServerError, "Failed to retrieve comments")
	}
	defer commentRows.Close()
	comments := []structs.Comment{}
	for commentRows.Next() {
		var comment structs.Comment
		err := commentRows.Scan(&comment.ID, &comment.Content, &comment.CreationDate, &comment.Username)
		if err != nil {
			return HttpError(forPage, http.StatusInternalServerError, "Failed to scan comment rows")
		}
		comments = append(comments, comment)
	}
	if err = commentRows.Err(); err != nil {
		return HttpError(forPage, http.StatusInternalServerError, "Failed to iterate over comment rows")
	}
	forPage.Posts = append(forPage.Posts, post)
	forPage.Comments = comments
	// Render the template with the data
	return forPage
}
