package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"gobang/app/api/global"
	"gobang/app/api/internal/model"
	"strconv"
)

type Game struct {
	RoomID     int `json:"roomID"`
	CheckBoard [][]int
}

func UpdateGame(r Game, roomID int) error {
	ctx := context.Background()
	roomJSON, err := json.Marshal(r)
	if err != nil {
		return err
	}
	err = global.RDB.HSet(ctx, "games", strconv.Itoa(roomID), roomJSON).Err()
	if err != nil {
		return err
	}
	return nil
}
func PutChess(c *gin.Context) {
	var m model.Message
	err := c.ShouldBindJSON(&m)
	if err != nil {
		global.Logger.Warn("get request failed," + err.Error())
		c.JSON(400, gin.H{
			"msg": "get request failed," + err.Error(),
		})
		return
	}
}
func StartGame(c *gin.Context) {
	var m model.Message
	err := c.ShouldBindJSON(&m)
	if err != nil {
		global.Logger.Warn("get request failed," + err.Error())
		c.JSON(400, gin.H{
			"msg": "get request failed," + err.Error(),
		})
		return
	}
	g := Game{
		RoomID:     m.RoomID,
		CheckBoard: make([][]int, 15),
	}
	for i := 0; i < 15; i++ {
		g.CheckBoard[i] = make([]int, 15)
	}
}
