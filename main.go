package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/haksunkim/flexiportal/app/controller"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Configuration struct {
	AppName string
	DB      struct {
		Host     string
		Name     string
		User     string
		Password string
		Port     int
	}
}

var Config Configuration

func main() {
	bac, _ := ioutil.ReadFile("config.json")
	err := json.Unmarshal(bac, &Config)
	if err != nil {
		log.Fatal(err)
	}

	dbConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", Config.DB.User, Config.DB.Password, Config.DB.Host, strconv.Itoa(Config.DB.Port), Config.DB.Name)

	controller.DBConn = dbConn

	r := mux.NewRouter()
	r.HandleFunc("/", controller.HomeHandler)
	r.PathPrefix("/resources/").Handler(http.StripPrefix("/resources/", http.FileServer(http.Dir("./resources"))))
	r.HandleFunc("/admin/main", controller.AdminMainHandler)
	r.HandleFunc("/admin/post/new", controller.NewPostHandler)
	r.HandleFunc("/admin/post", controller.CreatePostHandler).
		Methods("POST")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
