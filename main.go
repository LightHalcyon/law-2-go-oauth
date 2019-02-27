package main

import (
    "encoding/json"
    "log"
    "net/http"
	"github.com/gorilla/mux"
	"time"
	"os"
	"bytes"
)

// User struct
type User struct {
	UserID		int		`json:"userId,omitempty"`
	DisplayName	string	`json:"displayName,omitempty"`
}

// Comment struct
type Comment struct {
	ID			int	`json:"id,omitempty"`
	Comment		string	`json:"comment,omitempty"`
	CreatedBy	string	`json:"createdBy,omitempty"`
	CreatedAt	string	`json:"createdAt,omitempty"`
	UpdatedAt	string	`json:"updatedAt,omitempty"`
}

var users []User
var comments []Comment

// Login authentication to https://oauth.infralabs.cs.ui.ac.id/
func Login(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
}
func RegisterUser(w http.ResponseWriter, r *http.Request) {}
func GetUser(w http.ResponseWriter, r *http.Request) {}
func GetComment(w http.ResponseWriter, r *http.Request) {}
func GetCommentById(w http.ResponseWriter, r *http.Request) {}
func PostComment(w http.ResponseWriter, r *http.Request) {}
func DeleteComment(w http.ResponseWriter, r *http.Request) {}
func UpdateComment(w http.ResponseWriter, r *http.Request) {}

func main() {
	router := mux.NewRouter()
	// routes
	router.HandleFunc("/login", Login).Methods("POST")
	router.HandleFunc("/users", RegisterUser).Methods("POST")
	router.HandleFunc("/users", GetUser).Methods("GET")
	router.HandleFunc("/comments", GetComment).Methods("GET")
	router.HandleFunc("/comments", GetCommentById).Methods("GET")
	router.HandleFunc("/comments", PostComment).Methods("POST")
	router.HandleFunc("/comments", DeleteComment).Methods("HAPUS")
	router.HandleFunc("/comments", UpdateComment).Methods("UBAH")
	// execute
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", router))
}
