package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	// "os"
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
	Description string `json:"error_description"`
}

type OAuthData struct {
	AccessToken  	string
	UserID    		string
}

var users []User
var comments []Comment

// Authenticate validate access token given
func Authenticate(token string) (ErrorResponse, OAuthData){
	type oauthResponse struct {
		AccessToken  	string 	`json:"access_token"`
		Expires 		int 	`json:"expires"`
		UserID    		string 	`json:"user_id"`
		Scope        	string 	`json:"scope"`
		ClientID 		string 	`json:"client_id"`
	}

	var oauthresp OAuthData
	var oautherror ErrorResponse
	oauthURL := "https://oauth.infralabs.cs.ui.ac.id"
	verificationPath := "/oauth/resource"

	u, _ := url.ParseRequestURI(oauthURL)
	u.Path = verificationPath
	urlStr := u.String()

	client := &http.Client{}
	req, err := http.NewRequest("GET", urlStr, nil)
	req.Header.Add("Authorization", "Bearer " + token)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Request Error")
		errorresp := ErrorResponse{
			Status:      "error",
			Description: "500 Internal Server Error",
		}
		log.Println(err)
		return errorresp, oauthresp
	} else {
		defer resp.Body.Close()

		// var jsonstr jsonresp
		// var errorresp ErrorResponse
		dec := json.NewDecoder(resp.Body)
		if resp.StatusCode == http.StatusOK {
			decodeError := dec.Decode(&oauthresp)
			if decodeError != nil {
				log.Println("Decode Data Error")
				errorresp := ErrorResponse{
					Status:      "error",
					Description: "500 Internal Server Error",
				}
				log.Println(decodeError)
				return errorresp, oauthresp
			} else {
				if oauthresp.AccessToken == token {
					okReturn := OAuthData{
						AccessToken:	oauthresp.AccessToken,
						UserID:			oauthresp.UserID,
					}
					erroresp := ErrorResponse{
						Status:			"OK",
						Description:	"200 OK",
					}
					return erroresp, okReturn
				} else {
					errorresp := ErrorResponse{
						Status:			"error",
						Description:	"401 Unauthorized",
					}
					return errorresp, oauthresp
				}
			}
		} else {
			errorresp := ErrorResponse{
				Status:      "error",
				Description: resp.Status,
			}
			dec.Decode(&oautherror)
			log.Println(oautherror.Description)
			return errorresp, oauthresp
		}
	}
}

// Login authentication to https://oauth.infralabs.cs.ui.ac.id/
func Login(w http.ResponseWriter, r *http.Request) {
	type jsonresp struct {
		Status      string `json:"status"`
		AccessToken string `json:"token"`
	}

	type oauthResponse struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int 	`json:"expires_in"`
		TokenType    string `json:"token_type"`
		Scope        string `json:"scope"`
		RefreshToken string `json:"refresh_token"`
	}

	type receivedData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var params receivedData
	receiveddata := json.NewDecoder(r.Body)
	receiveddata.Decode(&params)
	log.Println(params)
	oauthURL := "https://oauth.infralabs.cs.ui.ac.id"
	tokenPath := "/oauth/token"
	// verificationPath := "/oauth/resource"

	data := url.Values{}
	data.Set("username", params.Username)
	data.Set("password", params.Password)
	data.Set("grant_type", "password")
	data.Set("client_id", "9c6xS7Z1XQWHzkLxMZHxvs0vmy0zFBUK")
	data.Set("client_secret", "hggpGtRNEMuU7nro4Z2WjODfB0Mdm3bc")
	log.Println(data)

	u, _ := url.ParseRequestURI(oauthURL)
	u.Path = tokenPath
	urlStr := u.String()

	client := &http.Client{}
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Request Error")
		errorresp := ErrorResponse{
			Status:      "error",
			Description: "500 Internal Server Error",
		}
		log.Println(err)
		json.NewEncoder(w).Encode(errorresp)
	} else {
		defer resp.Body.Close()

		var oautherror ErrorResponse
		var oauthresp oauthResponse
		// var jsonstr jsonresp
		// var errorresp ErrorResponse
		dec := json.NewDecoder(resp.Body)
		if resp.StatusCode == http.StatusOK {
			decodeError := dec.Decode(&oauthresp)
			if decodeError != nil {
				log.Println("Decode Data Error")
				errorresp := ErrorResponse{
					Status:      "error",
					Description: "500 Internal Server Error",
				}
				log.Println(decodeError)
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
			dec.Decode(&oautherror)
			log.Println(oautherror.Description)
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
