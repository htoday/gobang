package model

import (
	"github.com/dgrijalva/jwt-go"
	"sync"
)

type User struct {
	Username   string `form:"username" json:"username" db:"username"`
	Password   string `form:"password" json:"password" db:"password"`
	Nickname   string `form:"nickname" json:"nickname" db:"nickname"`
	Email      string `form:"email" json:"email" db:"email"`
	ID         int    `form:"uid" json:"uid" db:"id"`
	StarAmount int    `form:"starAmount" json:"starAmount" db:"starAmount"`
}
type MyClaims struct {
	Username string `form:"username" json:"username"`
	jwt.StandardClaims
}
type Message struct {
	Code      int    `form:"code" json:"code"`
	Username1 string `form:"username1" json:"username1"`
	Username2 string `form:"username2" json:"username2"`
	RoomID    string `form:"roomID" json:"roomID"`
	Row       int    `form:"row" json:"row"`
	Col       int    `form:"col" json:"col"`
}
type Counter struct {
	Mu    sync.Mutex
	Value int
}
