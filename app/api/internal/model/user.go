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
	Code         int    `form:"code" json:"code"`
	Username     string `form:"username" json:"username"`
	RoomID       int    `form:"roomID" json:"roomID"`
	Row          int    `form:"row" json:"row"`
	Col          int    `form:"col" json:"col"`
	RoomPassword string `form:"RoomPassword" json:"RoomPassword"`
	Ranking      bool   `form:"ranking" json:"ranking"`
	Forbidden    bool   `form:"forbidden" json:"forbidden"`
	Title        string `form:"title" json:"title"`
	FirstAct     int    `form:"firstAct" json:"firstAct"`
}
type Counter struct {
	Mu    sync.Mutex
	Value int
}
