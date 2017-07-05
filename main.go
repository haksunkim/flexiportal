package main

import (
	"github.com/gorilla/mux"
	"github.com/haksunkim/flexiportal/app/controller"
	"log"
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

type Configuration struct {
	AppName string
	DB struct {
		Host string
		User string
		Password string
		Port uint
	}
}

var Config Configuration

func main() {
	bac,_ := ioutil.ReadFile("config.json")
	json.Unmarshal(bac, &Config)
	
	r := mux.NewRouter()
	r.HandleFunc("/", controller.HomeHandler)
	r.PathPrefix("/resources/").Handler(http.StripPrefix("/resources/", http.FileServer(http.Dir("./resources"))))
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
