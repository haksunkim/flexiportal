package controller

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"time"
)

type Post struct {
	Id        uint
	Guid      string
	Title     string
	Content   string
	UserId    uint
	CreatedAt *time.Time
}
type Posts []Post

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseGlob("resources/templates/layout/*")
	t.ParseFiles("resources/templates/home/index.html")
	db, err := sql.Open("mysql", DBConn)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := db.Prepare("SELECT id, title, content, user_id, created_at FROM post;")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := stmt.Query()
	defer rows.Close()

	log.Printf("%s", rows)
	if err != nil {
		log.Fatal(err)
	}
	var posts Posts
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.UserId, &post.CreatedAt)

		if err != nil {
			log.Fatal(err)
		}
		posts = append(posts, post)
		log.Printf("%s: %s", post.Title, post.Content)
	}
	err = t.ExecuteTemplate(w, "layout", posts)
	if err != nil {
		log.Fatal(err)
	}
}
