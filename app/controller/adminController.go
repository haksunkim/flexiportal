package controller

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

var DBConn string

func AdminMainHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseGlob("resources/templates/layout/*")
	if err != nil {
		log.Fatal(err)
	}
	t.ParseFiles("resources/templates/admin/main.html")

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
	}
	err = t.ExecuteTemplate(w, "layout", posts)
	if err != nil {
		log.Fatal(err)
	}
}

func NewPostHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseGlob("resources/templates/layout/*")
	t.ParseFiles("resources/templates/admin/post/new.html")
	err := t.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		fmt.Printf(err.Error(), w)
	}
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	title := r.FormValue("title")
	content := r.FormValue("content")

	if len(title) > 0 && len(content) > 0 {
		db, err := sql.Open("mysql", DBConn)
		defer db.Close()
		if err != nil {
			log.Fatal(err)
		}
		stmt, err := db.Prepare("INSERT INTO post (title, content, user_id, created_at) VALUES (?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		_, err = stmt.Exec(title, content, 1, time.Now())
		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, r, "/admin/main", http.StatusSeeOther)
	} else {
		log.Printf("Title: %s", title)
	}
}
