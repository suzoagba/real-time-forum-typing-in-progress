package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"real-time-forum/internal/structs"
	"strings"
)

func PageHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageName := strings.Split(r.URL.Query().Get("page"), "=")

		pageData := structs.ForPage{}
		pageData.User = IsLoggedIn(r, db).User

		if pageData.User.LoggedIn {
			switch pageName[0] {
			case "logout":
				pageData = logout(w, db, pageData)
				jsonToPage(w, pageData)
				return
			case "createPost":
				pageData = createPost(r, db, pageData)
			case "viewPost":
				pageData = viewPost(db, pageData, pageName[1])
			case "comment":
				pageData = comment(r, db, pageData)
			default:
				pageData.Tags = structs.Tags
				var err error
				pageData.Posts, err = getAllPosts(db)
				if err != nil {
					pageData = HttpError(pageData, http.StatusInternalServerError, "Internal Server Error", "Error getting posts")
				}
			}
		} else {
			if r.Method == http.MethodPost {
				if pageName[0] == "login" {
					pageData = login(w, r, db, pageData)
				} else if pageName[0] == "register" {
					pageData = register(w, r, db, pageData)
				}
			}
		}

		jsonToPage(w, pageData)
		return
	}
}

func jsonToPage(w http.ResponseWriter, data structs.ForPage) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		HttpError(data, http.StatusInternalServerError, "Internal Server Error", "Error writing JSON")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
