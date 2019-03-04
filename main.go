package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

// User struct
type User struct {
	UserID      int    `json:"userId,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

// Comment struct
type Comment struct {
	ID        int    `json:"id,omitempty"`
	Comment   string `json:"comment,omitempty"`
	CreatedBy string `json:"createdBy,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
}

// ErrorResponse error response struct
type ErrorResponse struct {
	Status      string `json:"status"`
	Description string `json:"description"`
}

var users []User
var comments []Comment

// Login authentication to https://oauth.infralabs.cs.ui.ac.id/
func Login(w http.ResponseWriter, r *http.Request) {
	type jsonresp struct {
		Status      string `json:"status"`
		AccessToken string `json:"token"`
	}

	type oauthResponse struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    string `json:"expires_in"`
		TokenType    string `json:"token_type"`
		Scope        string `json:"scope"`
		RefreshToken string `json:"refresh_token"`
	}

	params := mux.Vars(r)
	oauthURL := os.Getenv("OAUTHURL")
	tokenPath := "/oauth/token"
	// verificationPath := "/oauth/resource"

	data := url.Values{}
	data.Set("username", params["username"])
	data.Set("password", params["password"])
	data.Set("grant_type", "password")
	data.Set("client_id", os.Getenv("CLIENTID"))
	data.Set("client_secret", os.Getenv("CLIENTSECRET"))

	u, _ := url.ParseRequestURI(oauthURL)
	u.Path = tokenPath
	urlStr := u.String()

	client := &http.Client{}
	r, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(r)
	if err != nil {
		errorresp := ErrorResponse{
			Status:      "error",
			Description: "500 Internal Server Error",
		}

		json.NewEncoder(w).Encode(errorresp)
	} else {
		defer resp.Body.Close()

		var oauthresp oauthResponse
		// var jsonstr jsonresp
		// var errorresp ErrorResponse
		dec := json.NewDecoder(resp.Body)
		if resp.StatusCode == http.StatusOK {
			decodeError := dec.Decode(&oauthresp)
			if decodeError != nil {
				errorresp := ErrorResponse{
					Status:      "error",
					Description: "500 Internal Server Error",
				}

				json.NewEncoder(w).Encode(errorresp)
			} else {
				jsonstr := jsonresp{
					Status:      resp.Status,
					AccessToken: oauthresp.AccessToken,
				}

				json.NewEncoder(w).Encode(jsonstr)
			}
		} else {
			errorresp := ErrorResponse{
				Status:      "error",
				Description: resp.Status,
			}

			json.NewEncoder(w).Encode(errorresp)
		}
	}
}
func RegisterUser(w http.ResponseWriter, r *http.Request)   {}
func GetUser(w http.ResponseWriter, r *http.Request)        {}
func GetComment(w http.ResponseWriter, r *http.Request)     {}
func GetCommentByID(w http.ResponseWriter, r *http.Request) {}
func PostComment(w http.ResponseWriter, r *http.Request)    {}
func DeleteComment(w http.ResponseWriter, r *http.Request)  {}
func UpdateComment(w http.ResponseWriter, r *http.Request)  {}

func main() {
	router := mux.NewRouter()
	// routes
	router.HandleFunc("/login", Login).Methods("POST")
	router.HandleFunc("/users", RegisterUser).Methods("POST")
	router.HandleFunc("/users", GetUser).Methods("GET")
	router.HandleFunc("/comments", GetComment).Methods("GET")
	router.HandleFunc("/comments", GetCommentByID).Methods("GET")
	router.HandleFunc("/comments", PostComment).Methods("POST")
	router.HandleFunc("/comments", DeleteComment).Methods("HAPUS")
	router.HandleFunc("/comments", UpdateComment).Methods("UBAH")
	// execute
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", router))
}
