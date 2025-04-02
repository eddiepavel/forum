package models

import (
	"html/template"
	"net/http"
	"strconv"
	"time"
)

var user = Users{
	ID:       1,
	Username: "Edouardos",
}

var comments = []Comment{
	{
		PostID:    1,
		Author:    user,
		Content:   template.HTML("This is a test comment<br>Great post!"),
		Upvotes:   2,
		Downvotes: 0,
		Time:      time.Now().Add(-time.Hour).Add(-time.Minute).Add(-time.Second).Format("2006-01-02 15:04:05"),
	},
	{
		PostID:    1,
		Author:    user,
		Content:   template.HTML("This is another test comment<br>Great post!"),
		Upvotes:   3,
		Downvotes: 5,
		Time:      time.Now().Add(-time.Hour).Add(-time.Minute).Add(-time.Second).Format("2006-01-02 15:04:05"),
	},
	{
		PostID:    2,
		Author:    user,
		Content:   template.HTML("Welcome!"),
		Upvotes:   0,
		Downvotes: 7,
		Time:      time.Now().Add(-time.Hour).Add(-time.Minute).Add(-time.Second).Format("2006-01-02 15:04:05"),
	},
}

var data = []Post{
	{
		ID:        1,
		Author:    user,
		Category:  "General",
		Content:   template.HTML("This is a test post<br>Hello, world!"),
		Title:     "Test Post",
		Time:      time.Now().Add(-time.Hour).Format("2006-01-02 15:04:05"),
		Upvotes:   5,
		Downvotes: 1,
	},
	{
		ID:       2,
		Author:   user,
		Category: "Announcements",
		Content:  "Welcome to the forum! Feel free to explore and contribute.",
		Title:    "Welcome Post",
		Time:     time.Now().Add(-2 * time.Hour).Format("2006-01-02 15:04:05"),
	},
}

func GetData(source string, r *http.Request) (interface{}, error) {
	// This function will be used to retrieve data from the database
	// and return it to the handlers
	if source == "home" {
		return data, nil
	} else if source == "view" {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			return nil, err
		}
		data[id-1].Comments = []Comment{}
		for _, comment := range comments {
			if comment.PostID == id {
				data[id-1].Comments = append(data[id-1].Comments, comment)
			}
		}
		return data[id-1], nil
	} else {
		return nil, nil
	}
}
