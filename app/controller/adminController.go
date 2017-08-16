package controller

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"database/sql"
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
	stmt, err := db.Prepare("SELECT id, title, content, user_id, created_at FROM blog;")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := stmt.Query()
	defer rows.Close()

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
	}
	err = t.ExecuteTemplate(w, "layout", blogs)
	if err != nil {
		log.Fatal(err)
	}
}

func NewBlogHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseGlob("resources/templates/layout/*")
	t.ParseFiles("resources/templates/admin/blog/new.html")
	err := t.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		fmt.Printf(err.Error(), w)
	}
}

func CreateBlogHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	title := r.FormValue("title")
	content := r.FormValue("content")

	if len(title) > 0 && len(content) > 0 {
		db, err := sql.Open("mysql", DBConn)
		defer db.Close()
		if err != nil {
			log.Fatal(err)
		}
		stmt, err := db.Prepare("INSERT INTO blog (title, content, user_id, created_at) VALUES (?, ?, ?, ?)")
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
