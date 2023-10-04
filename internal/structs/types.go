package structs

type Post struct {
	ID           int      `json:"id"`
	Username     string   `json:"username"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	CreationDate string   `json:"creationDate"`
	Tags         []string `json:"tags"`
}
type Comment struct {
	ID           string `json:"id"`
	Content      string `json:"content"`
	PostID       int    `json:"postID"`
	UserID       int    `json:"userID"`
	Username     string `json:"username"`
	CreationDate string `json:"creationDate"`
}
type User struct {
	ID       string `json:"userID"`
	Username string `json:"username"`
	LoggedIn bool   `json:"loggedIn"`
	Session  string `json:"session"`
}
type HttpError struct {
	Error bool   `json:"error"`
	Type  int    `json:"type"`
	Text  string `json:"text"`
	Text2 string `json:"text2"`
}
type ErrorMessage struct {
	Error   bool   `json:"error"`
	Type    bool   `json:"type"`
	Message string `json:"message"`
	Field1  string `json:"field1"`
	Field2  string `json:"field2"`
}
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type ForPage struct {
	HttpError HttpError    `json:"httpError"`
	Error     ErrorMessage `json:"error"`
	User      User         `json:"user"`
	Posts     []Post       `json:"posts"`
	Tags      []Tag        `json:"tags"`
	Comments  []Comment    `json:"comments"`
}
