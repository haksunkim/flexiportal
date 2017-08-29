package controller

import (
	"net/http"
	"log"
	"regexp"
	_ "github.com/lib/pq"
	"database/sql"
	"fmt"
	"crypto/sha1"
	"io"
	"github.com/gorilla/sessions"
	"github.com/haksunkim/flexiportal/app/entity"
	"time"
)

var sessionStore = sessions.NewCookieStore([]byte("flexiportal-application"))

func Register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	username := r.FormValue("username")
	email := r.FormValue("email")
	pass := r.FormValue("password")
	pageGUID := r.FormValue("referrer")
	// pass2 := r.FormValue("password2")
	gure := regexp.MustCompile("[^A-Za-z0-9]+")
	guid := gure.ReplaceAllString(username, "")
	password := weakPasswordHash(pass)

	db, err := sql.Open("postgres", DBConn)

	res, err := db.Exec("INSERT INTO user (username, guid, email, password) VALUES (?, ?, ?, ?)", username, guid, email, password)
	fmt.Println(res)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	} else {
		http.Redirect(w, r, "/page/" + pageGUID, 301)
	}
}

func weakPasswordHash(password string) []byte {
	hash := sha1.New()
	io.WriteString(hash, password)
	return hash.Sum(nil)
}

func getSessionUID(sid string) int {
	user := entity.User{}
	db, err := sql.Open("postgres", DBConn)

	err = db.QueryRow("SELECT user_id FROM session WHERE session_id=?", sid).Scan(user.Id)
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}

	return user.Id
}

func updateSession(sid string, uid int) {
	const timeFmt = "2006-01-02T15:04:05.999999999"
	tstamp := time.Now().Format(timeFmt)
	db, err := sql.Open("postgres", DBConn)
	_, err := db.Exec("INSERT INTO session (session_id, user_id, session_update) VALUES (?, ?, ?) ON CONFLICT UPDATE user_id=?, session_update=?", sid, uid, tstamp, uid, tstamp
	if err != nil {
		fmt.Println(err.Error())
	}
}