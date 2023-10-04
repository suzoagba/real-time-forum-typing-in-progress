package handlers

import (
	"database/sql"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client represents a connected WebSocket client.
type Client struct {
	conn *websocket.Conn
	User string
}
type Message struct {
	Type     string        `json:"type"`
	ID       int           `json:"id"`
	Message  string        `json:"content"`
	User     string        `json:"user"`
	To       string        `json:"to"`
	Users    []*OnlineUser `json:"users"`
	Time     string        `json:"time"`
	Messages []Message     `json:"messages"`
}
type OnlineUser struct {
	User        string       `json:"user"`
	Online      bool         `json:"online"`
	LastMessage *LastMessage `json:"lastMessage"`
}
type LastMessage struct {
	CreationDate time.Time `json:"time"`
	Receiver     string    `json:"to"`
	IsRead       bool      `json:"isRead"`
	Unread       int       `json:"unread"`
}

var clients = make(map[*Client]bool) // Map to store connected clients.
func WsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) { //, db *sql.DB) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	client := &Client{conn: conn, User: r.URL.Query().Get("username")}
	clients[client] = true
	// Notify all clients that a new user has joined.
	for c := range clients {
		broadcastMessage(Message{Type: "onlineUsers", User: c.User, Users: getClientUsernames(c.User, db)})
	}
	for {
		var message Message
		err := conn.ReadJSON(&message)
		if err != nil {
			if strings.Contains(err.Error(), "1000") || strings.Contains(err.Error(), "1001") || strings.Contains(err.Error(), "1005") {
				delete(clients, client)
				return
			} else {
				log.Printf("Error reading message: %v", err)
			}
		} else {
			if message.Type == "message" || message.Type == "getMessages" || message.Type == "read" {
				// Check if a chat exists between person1 and person2.
				chatID, err := getChatID(db, message.User, message.To)
				if err != nil {
					log.Fatal(err)
				}
				if chatID == 0 {
					// If no chat exists, create a new chat.
					chatID, err = createChat(db, message.User, message.To)
					if err != nil {
						log.Fatal(err)
					}
				}
				if message.Type == "message" {
					if err := addMessage(db, chatID, message.User, message.To, message.Message); err != nil {
						log.Fatal(err)
					}
				} else if message.Type == "getMessages" {
					messages, err := getMessageBatch(db, chatID, message.ID, 10)
					if err != nil {
						log.Fatal(err)
					}
					message.To = message.User
					if len(messages) != 0 {
						message.Messages = messages
					}
				} else if message.Type == "read" {
					err = markAsRead(db, chatID, message.To)
					if err != nil {
						log.Fatal(err)
					}
					continue
				}
			}
			broadcastMessage(message)
		}
	}
	// Notify all clients that a user has disconnected.
	for c := range clients {
		broadcastMessage(Message{Type: "onlineUsers", User: c.User, Users: getClientUsernames(c.User, db)})
	}
}
func getChatID(db *sql.DB, person1, person2 string) (int, error) {
	// Normalize the order of usernames to avoid duplicates.
	normalizedNames := []string{person1, person2}
	sqlQuery := "SELECT chatID FROM chats WHERE person1 = ? AND person2 = ? LIMIT 1;"
	if strings.Compare(person1, person2) > 0 {
		normalizedNames[0], normalizedNames[1] = normalizedNames[1], normalizedNames[0]
	}
	var chatID int
	err := db.QueryRow(sqlQuery, normalizedNames[0], normalizedNames[1]).Scan(&chatID)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return chatID, nil
}
func createChat(db *sql.DB, person1, person2 string) (int, error) {
	// Normalize the order of usernames to avoid duplicates.
	normalizedNames := []string{person1, person2}
	if strings.Compare(person1, person2) > 0 {
		normalizedNames[0], normalizedNames[1] = normalizedNames[1], normalizedNames[0]
	}
	result, err := db.Exec(`
        INSERT INTO chats (person1, person2)
        VALUES (?, ?);
    `, normalizedNames[0], normalizedNames[1])
	if err != nil {
		return 0, err
	}
	chatID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(chatID), nil
}
func addMessage(db *sql.DB, chatID int, sender, receiver, message string) error {
	_, err := db.Exec(`
        INSERT INTO messages (chatID, sender, receiver, message)
        VALUES (?, ?, ?, ?);
    `, chatID, sender, receiver, message)
	return err
}

// GetMessageBatch retrieves a batch of messages within a specific chat.
func getMessageBatch(db *sql.DB, chatID int, prevMessageID, batchSize int) ([]Message, error) {
	// SQL query to retrieve messages within a specific chat in batches.
	query := `
		SELECT messageID, sender, receiver, message, creationDate
		FROM messages
		WHERE chatID = ?
		ORDER BY messageID DESC
		LIMIT ?;
	`
	var args []interface{}
	if prevMessageID == 0 {
		args = []interface{}{chatID, batchSize}
	} else {
		query = strings.Replace(query, `WHERE chatID = ?`, `WHERE chatID = ? AND messageID < ?`, 1)
		args = []interface{}{chatID, prevMessageID, batchSize}
	}
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages []Message
	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.ID, &message.User, &message.To, &message.Message, &message.Time); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	if len(messages) == 0 {
		// No more messages available for this batch.
		return nil, nil
	}
	return messages, nil
}
func broadcastMessage(message Message) {
	for client := range clients {
		if (message.Type != "onlineUsers" && (client.User == message.To || client.User == message.User || message.To == "")) ||
			(message.Type == "onlineUsers" && client.User == message.User) {
			err := client.conn.WriteJSON(message)
			if err != nil {
				log.Printf("Error writing message: %v", err)
				client.conn.Close()
				delete(clients, client)
			}
		}
	}
}
func getClientUsernames(user string, db *sql.DB) []*OnlineUser {
	allUsers, err := getAllUsernames(db)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	usersList := make([]*OnlineUser, 0, len(allUsers))
	for clientTo := range allUsers {
		var err error
		chatId := 0
		to := allUsers[clientTo]
		if to != user {
			chatId, err = getChatID(db, user, to)
			if err != nil {
				log.Println(err)
			}
		}
		if chatId == 0 && to != user {
			usersList = append(usersList, &OnlineUser{User: to, Online: ifOnline(to)})
		} else {
			lastChat, err := getLastMessageByChatID(db, chatId)
			if err != nil {
				log.Println(err)
			}
			usersList = append(usersList, &OnlineUser{User: to, Online: ifOnline(to), LastMessage: lastChat})
		}
	}
	return usersList
}
func ifOnline(username string) bool {
	for client := range clients {
		if client.User == username {
			return true
		}
	}
	return false
}
func getAllUsernames(db *sql.DB) ([]string, error) {
	var usernames []string

	rows, err := db.Query("SELECT username FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, err
		}
		usernames = append(usernames, username)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return usernames, nil
}

func getLastMessageByChatID(db *sql.DB, chatID int) (*LastMessage, error) {
	query := `
		SELECT creationDate, receiver, isRead
		FROM messages
		WHERE chatID = ?
		ORDER BY creationDate DESC
		LIMIT 1;
	`
	var message LastMessage
	err := db.QueryRow(query, chatID).Scan(&message.CreationDate, &message.Receiver, &message.IsRead)
	if err != nil {
		if err == sql.ErrNoRows {
			// No messages found for the given chatID.
			return nil, nil
		}
		return nil, err
	}
	if message.IsRead == false {
		count, err := getUnreadMessageCount(db, chatID, message.Receiver)
		if err != nil {
			return nil, err
		}
		message.Unread = count
	}
	return &message, nil
}
func getUnreadMessageCount(db *sql.DB, chatID int, receiver string) (int, error) {
	// SQL query to count unread messages
	query := `
        SELECT COUNT(*) 
        FROM messages 
        WHERE chatID = ? AND receiver = ? AND isRead = 0;
    `

	var count int
	err := db.QueryRow(query, chatID, receiver).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
func markAsRead(db *sql.DB, chatID int, receiver string) error {
	query := `
        UPDATE messages
        SET isRead = true
        WHERE chatID = ? AND receiver = ?;
    `
	_, err := db.Exec(query, chatID, receiver)
	if err != nil {
		return err
	}
	return nil
}
