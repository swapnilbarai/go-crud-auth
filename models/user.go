package models

type User struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
	Email    string `json:"email"`
	MobileNo string `json:"mobileNo"`
}

// inmemory database for now
var Users map[string]User
