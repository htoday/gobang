package model

import "github.com/dgrijalva/jwt-go"

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}
type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
