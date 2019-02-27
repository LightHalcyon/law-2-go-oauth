package main

import (
    "encoding/json"
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	// routes
	router.HandleFunc("/login", Login).methods("POST")
	router.HandleFunc("/users", RegisterUser).methods("POST")
	router.HandleFunc("/users", GetUser).methods("GET")
	router.HandleFunc("/comments", GetComment).methods("POST")
	router.HandleFunc("/comments", DeleteComment).methods("DELETE")
	router.HandleFunc("/comments", UpdateComment).methods("PATCH")
	router.HandleFunc("/comments/add", PostComment).methods("POST")
	router.HandleFunc("/comments/{id}", GetCommentById).methods("GET")
	// execute
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", router))
}

