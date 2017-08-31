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
	"crypto/rand"
	"encoding/base64"
)

var UserSession entity.Session

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

func login(w http.ResponseWriter, r *http.Request) {
	validateSession(w, r)
	user := entity.User{}
	username := r.FormValue("username")
	pass := r.FormValue("password")
	password := weakPasswordHash(pass)
	db, err := sql.Open("postgres", DBConn)
	err = db.QueryRow("SELECT id, username FROM user WHERE username = ? and password = ?", username, password).Scan(&user.Id, &user.Name)
	if err != nil {
		fmt.Fprintln(w, err.Error())
		user.Id = 0
		user.Name = ""
	} else {
		updateSession(UserSession.Id, user.Id)
		fmt.Fprintln(w, user.Name)
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
	_, err = db.Exec("INSERT INTO session (session_id, user_id, session_update) VALUES (?, ?, ?) ON CONFLICT DO UPDATE SET user_id=?, session_update=?", sid, uid, tstamp, uid, tstamp)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func generateSessionId() string {
	sid := make([]byte, 24)
	_, err := io.ReadFull(rand.Reader, sid)
	if err != nil {
		log.Fatal("Could not generate session id")
	}
	return base64.URLEncoding.EncodeToString(sid)
}

func validateSession(w http.ResponseWriter, r *http.Request) {
	session, _ := sessionStore.Get(r, "app-session")
	if sid, valid := session.Values["sid"]; valid {
		currentUID := getSessionUID(sid.(string))
		updateSession(sid.(string), currentUID)
		UserSession.Id = string(currentUID)
	} else {
		newSID := generateSessionId()
		session.Values["sid"] = newSID
		session.Save(r, w)
		UserSession.Id = newSID
		updateSession(newSID, 0)
	}
	fmt.Println(session.ID)
}