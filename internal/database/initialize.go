package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func StartDB(input bool) {
	if input {
		err := createDatabase()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func createDatabase() error {
	// Create the database folder if it doesn't exist
	err := os.MkdirAll("./database", os.ModePerm)
	if err != nil {
		return err
	}
	// Create the database file
	_, err = os.Create("./database/forum.db")
	if err != nil {
		return err
	}
	return nil
}

func CreateTables(db *sql.DB) error {
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS users (
			userID TEXT PRIMARY KEY,
			username TEXT UNIQUE,
			email TEXT UNIQUE,
			age INTEGER,
			gender TEXT,
			firstName TEXT,
			lastName TEXT,
			password TEXT
		);
		CREATE TABLE IF NOT EXISTS posts (
			postID INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT,
			title TEXT,
			description TEXT,
			creationDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (username) REFERENCES users(username)
		);
		CREATE TABLE IF NOT EXISTS comments (
			commentID INTEGER PRIMARY KEY AUTOINCREMENT,
			postID INTEGER,
			username TEXT,
			content TEXT,
			creationDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (postID) REFERENCES posts(postID),
			FOREIGN KEY (username) REFERENCES users(username)
		);
		CREATE TABLE IF NOT EXISTS authenticated_users (
			sessionID TEXT PRIMARY KEY,
			username TEXT UNIQUE,
			FOREIGN KEY(username) REFERENCES users(username)
		);
		CREATE TABLE IF NOT EXISTS tags (
			tagID INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE
		);
		CREATE TABLE IF NOT EXISTS post_tags (
			postID INTEGER,
			tagID INTEGER,
			FOREIGN KEY (postID) REFERENCES posts(postID),
			FOREIGN KEY (tagID) REFERENCES tags(tagID),
			PRIMARY KEY (postID, tagID)
		);
		CREATE TABLE IF NOT EXISTS chats (
		  	chatID INTEGER PRIMARY KEY AUTOINCREMENT,
		  	person1 TEXT,
		  	person2 TEXT,
		  	FOREIGN KEY (person1) REFERENCES users(username),
		  	FOREIGN KEY (person2) REFERENCES users(username)
		);
		CREATE TABLE IF NOT EXISTS messages (
		    messageID INTEGER PRIMARY KEY AUTOINCREMENT,
		    chatID INTEGER,
		    sender TEXT,
		    receiver TEXT,
		    message TEXT,
		    creationDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		    isRead BOOLEAN DEFAULT 0,
		  	FOREIGN KEY (sender) REFERENCES users(username),
		  	FOREIGN KEY (receiver) REFERENCES users(username),
		  	FOREIGN KEY (chatID) REFERENCES chats(chatID)
		);
		INSERT OR IGNORE INTO tags (name) VALUES ('Cooking'), ('Mechanics'), ('Travel'), ('IT'), ('Random'), ('Market');`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		return err
	}
	fmt.Println("Database and tables created successfully!")
	return nil
}
