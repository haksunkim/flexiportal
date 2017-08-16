package controller

import (
	"html/template"
	"net/http"
	"database/sql"
	"time"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

type Blog struct {
	Id uint
	Title string
	Content string
	UserId uint
	CreatedAt *time.Time
}
type Blogs []Blog

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseGlob("resources/templates/layout/*")
	t.ParseFiles("resources/templates/home/index.html")
	db, err := sql.Open("mysql", DBConn)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := db.Prepare("SELECT id, title, content, user_id, created_at FROM blog;")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := stmt.Query()
	defer rows.Close()

	log.Printf("%s", rows)
	if err != nil {
		log.Fatal(err)
	}
	var blogs Blogs
	for rows.Next() {
		var blog Blog
		err := rows.Scan(&blog.Id, &blog.Title, &blog.Content, &blog.UserId, &blog.CreatedAt)

		if err != nil {
			log.Fatal(err)
		}
		blogs = append(blogs, blog)
		log.Printf("%s: %s", blog.Title, blog.Content)
	}
	err = t.ExecuteTemplate(w, "layout", blogs)
	if err != nil {
		log.Fatal(err)
	}
}
