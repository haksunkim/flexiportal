package main

import (
	"log"
	"net/http"
	"github.com/haksunkim/flexiportal/app/controller"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", controller.HomeHandler)
	http.Handle("/", r)
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
	}
	log.Fatal(srv.ListenAndServe())
}
