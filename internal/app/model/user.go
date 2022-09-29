package model

import "github.com/jinzhu/gorm"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type UserList struct {
	Users []DisplayUser `json:"users"`
}

type User struct {
	gorm.Model
	Email     string `gorm:"unique" json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type DisplayUser struct {
	Email     string `gorm:"unique" json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UserNameUpdate struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
