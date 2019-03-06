# LAW 2nd Assignment: Microservices with OAuth

This package starts an API Server tailored for 2nd Assignment of LAW Class in UI

## Install

With a [correctly configured](https://golang.org/doc/install#testing) Go toolchain:

```sh
go get github.com/reznov53/law-2-go-oauth
go install github.com/reznov53/law-2-go-oauth
```

The executable will be available at `$GOPATH/bin`

### Used environment variable

`$OAUTHURL`      = URL for OAuth server

`$CLIENTID`      = Client ID for OAuth server

`$CLIENTSECRET`  = Client Secret for OAuth server

## Available Routes

All routes below will return this payload if error

```json
{
    "status":       "error",
    "description":  "HTTPCODE Description"
}
```

### api/v1/login

#### POST

Logs the user in

Input Payload:

```json
{
    "username":   "username",
    "password":   "password"
}
```

Output Payload (200 OK):

```json
{
    "status":   "OK",
    "token":    "token-access-bearer"
}
```

### api/v1/users

#### POST

Register current logged in user with provided display name. Requires Authorization header with value `"Bearer access_token"`

Input Payload:

```json
{
	"displayName" : "Some User"
}
```

Output Payload:

```json
{
    "status" : "ok",
    "userId": 1,
    "displayName": "Some User" 
}
```

#### GET
Gets list of user based on parameters specified. Requires authenticated/logged in user with Authorization header with value `"Bearer access_token"`

Parameters:
```
page:   int
limit:  int
```

Output Payload:
```json
{
	"status" : "ok",
	"page" : 1,
	"limit" : 10,
	"total" : 2,
	"data" : [{
		"userId" : 1,
		"displayName" : "Some User"
	},{	
		"userId" : 2,
		"displayName" : "Another User"
	}]
}
```

### api/v1/comments/{id}
#### GET

Gets comment with provided id

```
{id}: int
```

Output Payload:
```json
{
	"status" : "ok",
	"data" : {
		"id" : 1,
		"comment" : "some comment",
		"createdBy" : "Some User",
		"createdAt" : "2018-02-13T08:34:57.000Z",
		"updatedAt" : "2018-02-13T08:34:57.000Z"
	}
}
```

### api/v1/comments
#### GET

Gets list of comments based on parameters specified

Parameters:
```
page:           int
limit:          int
createdBy:      String
startDate:      String (ISO 8601)
endDate:        String (ISO 8601)
```

Output Payload:
```json
{
	"status" : "ok",
	"page" : 1,
	"limit" : 10,
	"total" : 2,
	"data" : [{
		"id" : 1,
		"comment" : "some comment",
		"createdBy" : "Some User",
		"createdAt" : "2018-02-13T08:34:57.000Z",
		"updatedAt" : "2018-02-13T08:34:57.000Z"
	},{	
		"id" : 2,
		"comment" : "another comment",
		"createdBy" : "Another User",
		"createdAt" : "2018-02-14T01:20:21.000Z",
		"updatedAt" : "2018-02-14T01:20:21.000Z"
	}]
}
```

#### POST

Post a comment for the logged in user. Requires authenticated/logged in user with Authorization header with value `"Bearer access_token"`

Input Payload:
```json
{
	"comment" : "some comment"
}
```

Output Payload:
```json
{
	"status" : "ok",
	"data" : {
		"id" : 1,
		"comment" : "some comment",
		"createdBy" : "Some User",
		"createdAt" : "2018-02-13T08:34:57.000Z",
		"updatedAt" : "2018-02-13T08:34:57.000Z"
	}
}
```

### HAPUS

Deletes user's comment based on id. Requires authenticated/logged in user with Authorization header with value `"Bearer access_token"` & HTTP Method `HAPUS`

Input Payload:
```json
{ 
	"id" : 1 
}
```

Output Payload:
```json
{ 
	"status" : "ok" 
}
```

### UBAH

Updates user's comment with specified id. Requires authenticated/logged in user with Authorization header with value `"Bearer access_token"` & HTTP Method `UBAH`

Input Payload:
```json
{
	"id" : 2,
	"comment" : "updated comment"
}
```

Output Payload:
```json
{
	"status" : "ok",
	"data" : {
		"id" : 2,
		"comment" : "updated comment",
		"createdBy" : "Another User",
		"createdAt" : "2018-02-14T01:20:21.000Z",
		"updatedAt" : "2018-04-10T11:21:40.000Z"
	}
}
```