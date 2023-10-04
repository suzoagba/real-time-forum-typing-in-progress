package handlers

import (
	"database/sql"
	"real-time-forum/internal/structs"
	"strings"
)

func getAllPosts(db *sql.DB) ([]structs.Post, error) {
	query := `
		SELECT p.postID, u.username, p.title, p.description, p.creationDate, GROUP_CONCAT(t.name)
		FROM posts p
		JOIN users u ON p.username = u.username
		LEFT JOIN post_tags pt ON p.postID = pt.postID
		LEFT JOIN tags t ON pt.tagID = t.tagID
		GROUP BY p.postID
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts := []structs.Post{}
	for rows.Next() {
		var post structs.Post
		var tags string
		if err = rows.Scan(&post.ID, &post.Username, &post.Title, &post.Description, &post.CreationDate, &tags); err != nil {
			return nil, err
		}
		post.Tags = strings.Split(tags, ",")
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}
