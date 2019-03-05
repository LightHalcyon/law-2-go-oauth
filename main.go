package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"os"
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
var defaultTime time.Time

// IsExist checker method for users
func IsExist(users []User, user User) bool {
	for _, b := range users {
		if b.UserID == user.UserID || b.DisplayName == user.DisplayName {
			return true
		}
	}
	return false
}

// FilterCommentsByUName filters comment by user
func FilterCommentsByUName(comments []Comment, uname string) []Comment {
	var output []Comment
	for _, b := range comments {
		if b.CreatedBy == uname {
			output = append(output, b)
		}
	}
	return output
}

// FilterCommentsByDate filters comments by date
func FilterCommentsByDate(comments []Comment, startTime time.Time, endTime time.Time) []Comment {
	var output []Comment
	for _, b := range comments {
		updated, _ := time.Parse("2006-01-02T15:04:05-0700", b.UpdatedAt)
		if endTime != defaultTime {
			if updated.After(startTime) && updated.Before(endTime) {
				output = append(output, b)
			}
		} else {
			if updated.After(startTime) {
				output = append(output, b)
			}
		}
	}
	return output
}

// FindComment finds comment
func FindComment(comments []Comment, id int) (Comment, int) {
	var output Comment
	for i, b := range comments {
		if b.ID == id {
			return b, i
		}
	}
	return output, -1
}

// DelComment delete comment helper
func DelComment(comments []Comment, id int, userID string) ([]Comment, ErrorResponse) {
	err := ErrorResponse{
		Status:      "OK",
		Description: "200 OK",
	}
	uname := FindUName(users, userID)
	for i, b := range comments {
		if b.ID == id && b.CreatedBy == uname {
			comments = append(comments[:i], comments[i+1:]...)
			return comments, err
		} else if b.ID == id && b.CreatedBy != uname {
			err.Status = "error"
			err.Description = "401 Unauthorized"
			return comments, err
		}
	}
	return comments, err
}

// FindUName DisplayName finder
func FindUName(users []User, userID string) string {
	for _, b := range users {
		if b.ID == userID {
			return b.DisplayName
		}
	}
	return ""
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
	oauthURL := os.Getenv("OAUTHURL")
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
	oauthURL := os.Getenv("OAUTHURL")
	tokenPath := "/oauth/token"
	// verificationPath := "/oauth/resource"

	data := url.Values{}
	data.Set("username", params.Username)
	data.Set("password", params.Password)
	data.Set("grant_type", "password")
	data.Set("client_id", os.Getenv("CLIENTID"))
	data.Set("client_secret", os.Getenv("CLIENTSECRET"))
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
				// HeaderWriter(w, errorresp)
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
			// HeaderWriter(w, errorresp)
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
				// HeaderWriter(w, errorresp)
				json.NewEncoder(w).Encode(errorresp)
			}
		}
	} else {
		HeaderWriter(w, err)
		json.NewEncoder(w).Encode(err)
	}
}

// GetUser gets list of all user in users
func GetUser(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Status string `json:"status"`
		Page   int    `json:"page"`
		Limit  int    `json:"limit"`
		Total  int    `json:"total"`
		Data   []User `json:"data"`
	}

	authToken := strings.Split(r.Header.Get("Authorization"), " ")[1]
	err, _ := Authenticate(authToken)
	if err.Description == "200 OK" {
		params := r.URL.Query()
		page, _ := strconv.Atoi(params["page"][0])
		limit, _ := strconv.Atoi(params["limit"][0])
		var total int
		var data []User
		if len(users) <= limit || limit == 0 {
			data = users
			total = 1
		} else {
			total = len(users) / limit
			if (limit*page)-1 == 0 {
				data = users[0 : (limit*page)-1]
			} else {
				data = users[(limit*(page-1))-1 : (limit*page)-1]
			}
		}

		response := Response{
			Status: "OK",
			Page:   page,
			Limit:  limit,
			Total:  total,
			Data:   data,
		}
		json.NewEncoder(w).Encode(response)
	} else {
		HeaderWriter(w, err)
		json.NewEncoder(w).Encode(err)
	}
}

// PostComment post comment method
func PostComment(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Status string  `json:"status"`
		Data   Comment `json:"data"`
	}

	type Body struct {
		Comment string `json:"comment"`
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
			// HeaderWriter(w, errorresp)
			json.NewEncoder(w).Encode(errorresp)
		} else {
			displayName := FindUName(users, OAuthData.UserID)
			if displayName != "" {
				var id int
				if len(comments) == 0 {
					id = 1
				} else {
					id = comments[len(comments)-1].ID
				}
				comment := Comment{
					ID:        id,
					Comment:   body.Comment,
					CreatedBy: displayName,
					CreatedAt: time.Now().Format("2006-01-02T15:04:05-0700"),
					UpdatedAt: time.Now().Format("2006-01-02T15:04:05-0700"),
				}
				comments = append(comments, comment)
				response := Response{
					Status: "OK",
					Data:   comment,
				}
				json.NewEncoder(w).Encode(response)
			} else {
				errorresp := ErrorResponse{
					Status:      "error",
					Description: "401 Unauthorized",
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

// GetComment get comment list method
func GetComment(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Status string    `json:"status"`
		Page   int       `json:"page"`
		Limit  int       `json:"limit"`
		Total  int       `json:"total"`
		Data   []Comment `json:"data"`
	}
	
	params := r.URL.Query()
	page, _ := strconv.Atoi(params["page"][0])
	limit, _ := strconv.Atoi(params["limit"][0])
	createdBy := params["createdBy"][0]
	var startDate time.Time
	var endDate time.Time
	if params["startDate"][0] != "" {
		startDate, _ = time.Parse("2006-01-02T15:04:05-0700", params["startDate"][0])
	}
	if params["endDate"][0] != "" {
		endDate, _ = time.Parse("2006-01-02T15:04:05-0700", params["endDate"][0])
	}
	var total int
	data := comments
	if createdBy != "" {
		data = FilterCommentsByUName(data, createdBy)
	}
	data = FilterCommentsByDate(data, startDate, endDate)
	if len(data) <= limit || limit == 0 {
		total = 1
	} else {
		total = len(data) / limit
		if (limit*page)-1 == 0 {
			data = data[0 : (limit*page)-1]
		} else {
			data = data[(limit*(page-1))-1 : (limit*page)-1]
		}
	}

	response := Response{
		Status: "OK",
		Page:   page,
		Limit:  limit,
		Total:  total,
		Data:   data,
	}
	json.NewEncoder(w).Encode(response)
}

// GetCommentByID search comments for comment with id = {id}
func GetCommentByID(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Status string  `json:"status"`
		Data   Comment `json:"data"`
	}
	ID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		errorresp := ErrorResponse{
			Status:      "error",
			Description: "500 Internal Server Error",
		}
		// HeaderWriter(w, errorresp)
		json.NewEncoder(w).Encode(errorresp)
	} else {
		comment, _ := FindComment(comments, ID)
		response := Response{
			Status: "OK",
			Data:   comment,
		}
		json.NewEncoder(w).Encode(response)
	}
}

// DeleteComment delete comment with id
func DeleteComment(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Status string `json:"status"`
	}

	type Body struct {
		ID int `json:"id"`
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
			// HeaderWriter(w, errorresp)
			json.NewEncoder(w).Encode(errorresp)
		} else {
			comment, errors := DelComment(comments, body.ID, OAuthData.UserID)
			if errors.Description == "200 OK" {
				comments = comment
				response := Response{
					Status: "OK",
				}
				json.NewEncoder(w).Encode(response)
			} else {
				HeaderWriter(w, errors)
				json.NewEncoder(w).Encode(errors)
			}
		}
	} else {
		HeaderWriter(w, err)
		json.NewEncoder(w).Encode(err)
	}
}

// UpdateComment updates comment based on it's id
func UpdateComment(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Status string  `json:"status"`
		Data   Comment `json:"data"`
	}

	type Body struct {
		ID      int    `json:"id"`
		Comment string `json:"comment"`
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
			// HeaderWriter(w, errorresp)
			json.NewEncoder(w).Encode(errorresp)
		} else {
			commentItem, i := FindComment(comments, body.ID)
			if commentItem.CreatedBy != FindUName(users, OAuthData.UserID) {
				errorresp := ErrorResponse{
					Status:      "error",
					Description: "401 Unauthorized",
				}
				HeaderWriter(w, errorresp)
				json.NewEncoder(w).Encode(errorresp)
			} else {
				commentItem.Comment = body.Comment
				commentItem.UpdatedAt = time.Now().Format("2006-01-02T15:04:05-0700")
				comments[i] = commentItem
				response := Response{
					Status: "OK",
					Data:   commentItem,
				}
				json.NewEncoder(w).Encode(response)
			}
		}
	} else {
		HeaderWriter(w, err)
		json.NewEncoder(w).Encode(err)
	}
}

func main() {
	router := mux.NewRouter()
	// routes
	router.HandleFunc("/api/v1/login", Login).Methods("POST")
	router.HandleFunc("/api/v1/users", RegisterUser).Methods("POST")
	router.HandleFunc("/api/v1/users", GetUser).Methods("GET")
	router.HandleFunc("/api/v1/comments", GetComment).Methods("GET")
	router.HandleFunc("/api/v1/comments/{id}", GetCommentByID).Methods("GET")
	router.HandleFunc("/api/v1/comments", PostComment).Methods("POST")
	router.HandleFunc("/api/v1/comments", DeleteComment).Methods("HAPUS")
	router.HandleFunc("/api/v1/comments", UpdateComment).Methods("UBAH")
	// execute
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", router))
}
