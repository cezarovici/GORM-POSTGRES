package models

import (
	jwt "github.com/golang-jwt/jwt/v5"
)

// User represents a User schema
type User struct {
	Base
	Email    string `json:"email" gorm:"unique"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
}

type UserErrors struct {
	Err      bool   `json:"erors"`
	Email    string `json: "email"`
	Username string `json: "username"`
	Password string `json: "password"`
}

type Claims struct {
	jwt.RegisteredClaims
	ID uint `gorm: "primaryKey"`
}
