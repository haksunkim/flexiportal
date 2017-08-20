package controller

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
	"strings"
	"github.com/lib/pq"
)

var DBConn string

func AdminMainHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseGlob("resources/templates/layout/*")
	if err != nil {
		log.Fatal(err)
	}
	t.ParseFiles("resources/templates/admin/main.html")

	db, err := sql.Open("postgres", DBConn)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := db.Prepare("SELECT id, title, content, created_by, created_at FROM post;")
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
		err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.CreatedBy, &post.CreatedAt)

		if err != nil {
			log.Fatal(err)
		}
		post.StrCreatedAt = post.CreatedAt.Format("02/01/2006 15:04:05PM")
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

	guid := strings.Replace(title, " ", "-", -1)

	if len(title) > 0 && len(content) > 0 {
		db, err := sql.Open("postgres", DBConn)
		defer db.Close()
		if err != nil {
			log.Fatal(err)
		}
		_, err = db.Exec("INSERT INTO post (guid, title, content, created_by, created_at) VALUES ($1, $2, $3, $4, $5)", guid, title, content, 1, pq.FormatTimestamp(time.Now()))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Redirect(w, r, "/admin/main", http.StatusSeeOther)
		}
	} else {
		log.Printf("Title: %s", title)
	}
}
