package main

import (
	"strconv"
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
	ID          string `json:"-"`
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

// OAuthData response struct
type OAuthData struct {
	AccessToken string `json:"access_token,omitempty"`
	ClientID    string `json:"client_id,omitempty"`
	UserID      string `json:"user_id,omitempty"`
}

var users []User
var comments []Comment

// IsExist checker method for users
func IsExist(users []User, user User) bool {
	for _, b := range users {
		if b.UserID == user.UserID || b.DisplayName == user.DisplayName {
			return true
		}
	}
	return false
}

// HeaderWriter helper method for response returns
func HeaderWriter(w http.ResponseWriter, err ErrorResponse) {
	if err.Description == "500 Internal Server Error" {
		w.WriteHeader(http.StatusInternalServerError)
	} else if err.Description == "401 Unauthorized" {
		w.WriteHeader(http.StatusUnauthorized)
	} else if err.Description == "200 OK" {
		w.WriteHeader(http.StatusOK)
	} else if err.Description == "403 Forbidden" {
		w.WriteHeader(http.StatusForbidden)
	}
}

// Authenticate validate access token given
func Authenticate(token string) (ErrorResponse, OAuthData) {
	var oauthresp OAuthData
	var oautherror ErrorResponse
	var errorresp ErrorResponse
	oauthURL := "https://oauth.infralabs.cs.ui.ac.id"
	verificationPath := "/oauth/resource"

	u, _ := url.ParseRequestURI(oauthURL)
	u.Path = verificationPath
	urlStr := u.String()

	client := &http.Client{}
	req, err := http.NewRequest("GET", urlStr, nil)
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Request Error")
		errorresp = ErrorResponse{
			Status:      "error",
			Description: "500 Internal Server Error",
		}
		log.Println(err)
	} else {
		defer resp.Body.Close()
		log.Println(resp.Body)
		dec := json.NewDecoder(resp.Body)
		log.Println(resp.Status)
		if resp.StatusCode == http.StatusOK {
			decodeError := dec.Decode(&oauthresp)
			log.Println(oauthresp)
			if decodeError != nil {
				log.Println("Decode Data Error")
				errorresp = ErrorResponse{
					Status:      "error",
					Description: "500 Internal Server Error",
				}
			} else {
				if oauthresp.AccessToken == token {
					oauthresp = OAuthData{
						AccessToken: oauthresp.AccessToken,
						UserID:      oauthresp.UserID,
					}
					errorresp = ErrorResponse{
						Status:      "OK",
						Description: "200 OK",
					}
				} else {
					errorresp = ErrorResponse{
						Status:      "error",
						Description: "401 Unauthorized",
					}
				}
			}
		} else {
			errorresp = ErrorResponse{
				Status:      "error",
				Description: resp.Status,
			}
			dec.Decode(&oautherror)
			log.Println(oautherror.Description)
		}
	}
	return errorresp, oauthresp
}

// Login authentication to https://oauth.infralabs.cs.ui.ac.id/
func Login(w http.ResponseWriter, r *http.Request) {
	type jsonresp struct {
		Status      string `json:"status"`
		AccessToken string `json:"token"`
	}

	type oauthResponse struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
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
				HeaderWriter(w, errorresp)
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
			HeaderWriter(w, errorresp)
			json.NewEncoder(w).Encode(errorresp)
			dec.Decode(&oautherror)
			log.Println(oautherror.Description)
		}
	}
}

// RegisterUser register user to db
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	// var errorResp ErrorResponse
	type Body struct {
		DisplayName string `json:"displayName"`
	}

	type Response struct {
		Status      string `json:"status"`
		UserID      int    `json:"userId"`
		DisplayName string `json:"displayName"`
	}

	authToken := strings.Split(r.Header.Get("Authorization"), " ")[1]
	err, OAuthData := Authenticate(authToken)
	if err.Description == "200 OK" {
		var body Body
		data := json.NewDecoder(r.Body)
		error := data.Decode(&body)
		if error != nil {
			errorresp := ErrorResponse{
				Status:      "error",
				Description: "500 Internal Server Error",
			}
			HeaderWriter(w, errorresp)
			json.NewEncoder(w).Encode(errorresp)
		} else {
			user := User{
				UserID:      len(users) + 1,
				DisplayName: body.DisplayName,
				ID:          OAuthData.UserID,
			}
			if !IsExist(users, user) {
				users = append(users, user)
				response := Response{
					Status:      "OK",
					UserID:      user.UserID,
					DisplayName: user.DisplayName,
				}
				json.NewEncoder(w).Encode(response)
			} else {
				errorresp := ErrorResponse{
					Status:      "error",
					Description: "403 Forbidden",
				}
				HeaderWriter(w, errorresp)
				json.NewEncoder(w).Encode(errorresp)
			}
		}
	} else {
		HeaderWriter(w, err)
		json.NewEncoder(w).Encode(err)
	}
}

// GetUser gets list of all user in users
func GetUser(w http.ResponseWriter, r *http.Request)        {
	type Response struct {
		Status	string	`json:"status"`
		Page	int		`json:"page"`
		Limit	int		`json:"limit"`
		Total	int		`json:"total"`
		Data	[]User	`json:"data"`
	}

	authToken := strings.Split(r.Header.Get("Authorization"), " ")[1]
	err, _ := Authenticate(authToken)
	if err.Description == "200 OK" {
		params := r.URL.Query()
		log.Println(params)
		page, _ := strconv.Atoi(params["page"][0])
		limit, _ := strconv.Atoi(params["limit"][0])
		var total int
		var data []User
		if len(users) <= limit || limit == 0 {
			data = users
			total = 1
		} else {
			total = len(users)/limit
			if (limit*page)-1 == 0 {
				data = users[0:(limit*page)-1]
			} else {
				data = users[(limit*(page-1))-1:(limit*page)-1]
			}
		}
	
		response := Response{
			Status:	"OK",
			Page:	page,
			Limit:	limit,
			Total:	total,
			Data:	data,
		}
		json.NewEncoder(w).Encode(response)
	} else {
		HeaderWriter(w, err)
		json.NewEncoder(w).Encode(err)
	}
}

// PostComment post comment method
func PostComment(w http.ResponseWriter, r *http.Request)    {

}
func GetComment(w http.ResponseWriter, r *http.Request)     {}
func GetCommentByID(w http.ResponseWriter, r *http.Request) {}
func DeleteComment(w http.ResponseWriter, r *http.Request)  {}
func UpdateComment(w http.ResponseWriter, r *http.Request)  {}

func main() {
	router := mux.NewRouter()
	// routes
	router.HandleFunc("/api/v1/login", Login).Methods("POST")
	router.HandleFunc("/api/v1/users", RegisterUser).Methods("POST")
	router.HandleFunc("/api/v1/users", GetUser).Methods("GET")
	router.HandleFunc("/api/v1/comments", GetComment).Methods("GET")
	router.HandleFunc("/api/v1/comments", GetCommentByID).Methods("GET")
	router.HandleFunc("/api/v1/comments", PostComment).Methods("POST")
	router.HandleFunc("/api/v1/comments", DeleteComment).Methods("HAPUS")
	router.HandleFunc("/api/v1/comments", UpdateComment).Methods("UBAH")
	// execute
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", router))
}
