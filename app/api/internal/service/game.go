package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"gobang/app/api/global"
	"gobang/app/api/internal/dao/mysql"
	"gobang/app/api/internal/model"
	"math/rand"
	"strconv"
	"time"
)

type Game struct {
	RoomID     int     `json:"roomID"`
	CheckBoard [][]int `json:"gameBoard"`
	Turn       int     `json:"turn"`
	FirstAct   int     `json:"firstAct"`
}

func GetGame(roomID int) (Game, error) {
	var thisGame Game
	ctx := context.Background()
	roomJSON, err := global.RDB.HGet(ctx, "games", strconv.Itoa(roomID)).Result()
	if err != nil {
		global.Logger.Warn("translate JSON into game failed," + err.Error())
		return Game{}, err
	}
	//反序列化
	err = json.Unmarshal([]byte(roomJSON), &thisGame)
	if err != nil {
		global.Logger.Warn("Unmarshal JSON failed," + err.Error())
		return thisGame, err
	}
	return thisGame, nil
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
	thisRoom, err := GetRoom(m.RoomID)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "get room JSON failed," + err.Error(),
		})
		return
	}
	var thisGame Game
	thisGame, err = GetGame(m.RoomID)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "get game JSON failed," + err.Error(),
		})
		return
	}

	//以下为下棋
	if thisGame.Turn == 1 && m.Username != thisRoom.User1 {
		c.JSON(400, gin.H{
			"msg": "It is not your turn," + err.Error(),
		})
		return
	}
	if thisGame.Turn == 2 && m.Username != thisRoom.User2 {
		c.JSON(400, gin.H{
			"msg": "It is not your turn," + err.Error(),
		})
		return
	}
	if thisGame.FirstAct == 1 { //玩家一执黑
		if thisGame.Turn == 1 { //轮到玩家一行动
			if thisRoom.Forbidden == true {
				flag := JudgeBan(thisGame.CheckBoard, m.Row, m.Col)
				if flag == true {
					c.JSON(400, gin.H{
						"msg": "Forbidden point",
					})
					return
				}
			}
			thisGame.CheckBoard[m.Row][m.Col] = 1 //1为黑子
			if CheckWin(thisGame.CheckBoard, 1, m.Row, m.Col) {
				thisRoom.RoomStatus = 3
				if thisRoom.Ranking == true {
					mysql.UpdateStar(thisRoom.User1, thisRoom.User2)
				}
				err = UpdateRoom(thisRoom, thisRoom.RoomID)
				if err != nil {
					global.Logger.Warn("Update Room Ready failed," + err.Error())
					c.JSON(400, gin.H{
						"msg": "Update Room Ready failed," + err.Error(),
					})
					return
				}
				go ResetRoom(thisGame.RoomID)
			}
		} else { //轮到玩家二行动
			thisGame.CheckBoard[m.Row][m.Col] = 2 //2为白子
			if CheckWin(thisGame.CheckBoard, 2, m.Row, m.Col) {
				thisRoom.RoomStatus = 4
				if thisRoom.Ranking == true {
					mysql.UpdateStar(thisRoom.User2, thisRoom.User1)
				}
				err = UpdateRoom(thisRoom, thisRoom.RoomID)
				if err != nil {
					global.Logger.Warn("Update Room Ready failed," + err.Error())
					c.JSON(400, gin.H{
						"msg": "Update Room Ready failed," + err.Error(),
					})
					return
				}
				go ResetRoom(thisGame.RoomID)
			}

		}
	} else { //玩家2执黑
		if thisGame.Turn == 1 { //轮到玩家一行动
			thisGame.CheckBoard[m.Row][m.Col] = 2
			if CheckWin(thisGame.CheckBoard, 2, m.Row, m.Col) {
				if thisRoom.Ranking == true {
					mysql.UpdateStar(thisRoom.User1, thisRoom.User2)
				}
				thisRoom.RoomStatus = 3
				err = UpdateRoom(thisRoom, thisRoom.RoomID)
				if err != nil {
					global.Logger.Warn("Update Room Ready failed," + err.Error())
					c.JSON(400, gin.H{
						"msg": "Update Room Ready failed," + err.Error(),
					})
					return
				}
				go ResetRoom(thisGame.RoomID)
			}
		} else { //轮到玩家二行动

			if thisRoom.Forbidden == true {
				flag := JudgeBan(thisGame.CheckBoard, m.Row, m.Col)
				if flag == true {
					c.JSON(400, gin.H{
						"msg": "Forbidden point",
					})
					return
				}
			}
			thisGame.CheckBoard[m.Row][m.Col] = 1 //1为黑子
			if CheckWin(thisGame.CheckBoard, 1, m.Row, m.Col) {
				thisRoom.RoomStatus = 4
				if thisRoom.Ranking == true {
					mysql.UpdateStar(thisRoom.User2, thisRoom.User1)
				}
				err = UpdateRoom(thisRoom, thisRoom.RoomID)
				if err != nil {
					global.Logger.Warn("Update Room Ready failed," + err.Error())
					c.JSON(400, gin.H{
						"msg": "Update Room Ready failed," + err.Error(),
					})
					return
				}
				go ResetRoom(thisGame.RoomID)
			}
		}
	}
	if thisGame.Turn == 1 {
		thisGame.Turn = 2
	} else {
		thisGame.Turn = 1
	}

	err = UpdateGame(thisGame, thisRoom.RoomID)
	if err != nil {
		global.Logger.Warn("save game information failed," + err.Error())
		c.JSON(400, gin.H{
			"msg": "save game information failed," + err.Error(),
		})
		return
	}
	global.Logger.Info("put chess success")
	c.JSON(200, gin.H{
		"msg": "put chess success",
	})
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

	var thisRoom Room
	thisRoom, err = GetRoom(m.RoomID)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "get room JSON failed," + err.Error(),
		})
		return
	}
	if thisRoom.User1Ready != true || thisRoom.User2Ready != true {
		global.Logger.Warn("User not ready")
		c.JSON(400, gin.H{
			"msg": "User not ready",
		})
		return
	}
	thisRoom.RoomStatus = 1
	g := Game{
		RoomID:     m.RoomID,
		CheckBoard: make([][]int, 15),
		Turn:       thisRoom.FirstAct,
	}
	for i := 0; i < 15; i++ {
		g.CheckBoard[i] = make([]int, 15)
	}
	if thisRoom.FirstAct == 3 {
		//设置当前时间为种子
		rand.Seed(time.Now().UnixNano())
		// 生成随机整数在[1, 2]范围内
		randNum := rand.Intn(2) + 1
		g.Turn = randNum
		g.FirstAct = randNum
	} else {
		g.FirstAct = thisRoom.FirstAct
	}
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			g.CheckBoard[i][j] = 0
		}
	}
	err = UpdateRoom(thisRoom, thisRoom.RoomID)
	if err != nil {
		global.Logger.Warn("save room information failed," + err.Error())
		c.JSON(400, gin.H{
			"msg": "save room information failed," + err.Error(),
		})
		return
	}
	err = UpdateGame(g, thisRoom.RoomID)
	if err != nil {
		global.Logger.Warn("save game information failed," + err.Error())
		c.JSON(400, gin.H{
			"msg": "save game information failed," + err.Error(),
		})
		return
	}
	global.Logger.Info("start game!! in room " + strconv.Itoa(thisRoom.RoomID))
	c.JSON(200, gin.H{
		"msg": "start game!!",
	})
}
func JudgeBan(a [][]int, x int, y int) bool {
	a[x][y] = 1 //1为黑子
	//先判断胜利，再长连禁手
	flag := CheckWin(a, 1, x, y)
	if flag == true {
		//长连这里被我写成了检查是否存在5连，如果有5连，则不是禁手，反之长连
		flag = JudgeLongForbidden(a, x, y)
		if flag == true {
			a[x][y] = 0
			return false
		} else {
			a[x][y] = 0
			return true
		}
	}
	flag = JudgeFourFourForbidden(a, x, y)
	if flag == true {
		a[x][y] = 0
		return true
	}
	flag = JudgeThreeThreeForbidden(a, x, y)
	if flag == true {
		a[x][y] = 0
		return true
	}
	a[x][y] = 0
	return false
}
func JudgeFourFourForbidden(a [][]int, x int, y int) bool {
	fourCount := getFourOneLine(a, x, y, 1, 0) + getFourOneLine(a, x, y, 0, 1) + getFourOneLine(a, x, y, -1, 1) + getFourOneLine(a, x, y, 1, 1)
	return fourCount > 1
}
func JudgeThreeThreeForbidden(a [][]int, x int, y int) bool {
	ThreeCount := getThreeOneLine(a, x, y, 1, 0) + getThreeOneLine(a, x, y, 0, 1) + getThreeOneLine(a, x, y, -1, 1) + getThreeOneLine(a, x, y, 1, 1)
	return ThreeCount > 1
}
func getThreeOneLine(a [][]int, x int, y int, dx int, dy int) int {
	beforeBefore := 0
	before := 0
	middle := 1
	after := 0
	afterAfter := 0
	flag := 0
	i := x - dx
	j := y - dy
	for {
		if !(i >= 0 && i <= 14 && j >= 0 && j <= 14) || a[i][j] == 2 {
			if flag == 0 {
				before = -1
			}
			if flag == 1 {
				beforeBefore = -1
			}
			break
		}
		if a[i][j] == 1 {
			if flag == 1 {
				before++
			} else if flag == 2 {
				beforeBefore++
			} else {
				middle++
			}
		} else if flag == 2 {
			break
		} else {
			flag++
		}
		i -= dx
		j -= dy
	}
	deltaMiddle := middle
	flag = 0
	i = x + dx
	j = y + dy
	for {
		if !(i >= 0 && i <= 14 && j >= 0 && j <= 14) || a[i][j] == 2 {
			if flag == 0 {
				after = -1
			}
			if flag == 1 {
				afterAfter = -1
			}
			break
		}
		if a[i][j] == 1 {
			if flag == 1 {
				after++
			} else if flag == 2 {
				afterAfter++
			} else {
				middle++
			}
		} else if flag == 2 {
			break
		} else {
			flag++
		}
		i += dx
		j += dy
	}

	//情况判断
	if before == -1 || after == -1 || middle > 3 {
		return 0
	}
	if middle == 3 {
		if before != 0 || after != 0 {
			return 0
		}
		if beforeBefore == -1 && afterAfter == -1 {
			return 0
		}
		if beforeBefore == -1 {
			if afterAfter > 0 {
				return 0
			}
			if JudgeBan(a, x+dx*(4-deltaMiddle), y+dy*(4-deltaMiddle)) {
				return 0
			} else {
				return 1
			}

		}
		if afterAfter == -1 {
			if beforeBefore > 0 {
				return 0
			}
			if JudgeBan(a, x-dx*deltaMiddle, y-dy*deltaMiddle) {
				return 0
			} else {
				return 1
			}
		}
		if beforeBefore > 0 && afterAfter > 0 {
			return 0
		}
		a1 := JudgeBan(a, x-dx*deltaMiddle, y-dy*deltaMiddle)
		b1 := JudgeBan(a, x+dx*(4-deltaMiddle), y+dy*(4-deltaMiddle))
		if a1 && b1 {
			return 0
		} else {
			return 1
		}
	}
	if middle == 2 {
		if before == 1 && after == 0 {
			if beforeBefore == -1 {
				return 0
			}
			if JudgeBan(a, x-dx*deltaMiddle, y-dy*deltaMiddle) {
				return 0
			} else {
				return 1
			}
		}
		if before == 0 && after == 1 {
			if afterAfter == -1 {
				return 0
			}
			if JudgeBan(a, x+dx*(3-deltaMiddle), y+dy*(3-deltaMiddle)) {
				return 0
			} else {
				return 1
			}
		}
		return 0
	}
	if middle == 1 {
		if before == 2 && after == 0 {
			if beforeBefore == -1 {
				return 0
			}
			if JudgeBan(a, x-dx, y-dy) {
				return 0
			} else {
				return 1
			}
		}
		if before == 0 && after == 2 {
			if afterAfter == -1 {
				return 0
			}
			if JudgeBan(a, x+dx, y+dy) {
				return 0
			} else {
				return 1
			}
		}

	}
	return 0
}
func getFourOneLine(a [][]int, x int, y int, dx int, dy int) int {
	before := 0
	middle := 1
	after := 0
	flag := false
	i := x - dx
	j := y - dy
	for {
		if !(i >= 0 && i <= 14 && j >= 0 && j <= 14) || a[i][j] == 2 {
			if !flag {
				before = -1
			}
			break
		}
		if a[i][j] == 1 {
			if flag {
				before++
			} else {
				middle++
			}
		} else {
			if flag {
				break
			}
			flag = true
		}
		i -= dx
		j -= dy
	}
	flag = false
	i = x + dx
	j = y + dy
	for {
		if !(i >= 0 && i <= 14 && j >= 0 && j <= 14) || a[i][j] == 2 {
			if !flag {
				after = -1
			}
			break
		}
		if a[i][j] == 1 {
			if flag {
				after++
			} else {
				middle++
			}
		} else {
			if flag {
				break
			}
			flag = true
		}
		i += dx
		j += dy
	}
	if middle == 4 {
		if before == 0 || after == 0 {
			return 1
		}
		return 0
	}
	var sum int
	sum = 0
	if before > 0 && before+middle == 4 {
		sum++
	}
	if after > 0 && after+middle == 4 {
		sum++
	}
	return sum
}
func JudgeLongForbidden(a [][]int, x int, y int) bool {
	var length int
	color := 1
	//横着找
	length = 1
	for i := x - 1; i >= 0; i-- {
		if a[i][y] == color {
			length++
		} else {
			break
		}
	}
	for i := x + 1; i <= 14; i++ {
		if a[i][y] == color {
			length++
		} else {
			break
		}
	}
	if length == 5 {
		return true
	}
	//竖着找
	length = 1
	for i := y - 1; i >= 0; i-- {
		if a[x][i] == color {
			length++
		} else {
			break
		}
	}
	for i := y + 1; i <= 14; i++ {
		if a[x][i] == color {
			length++
		} else {
			break
		}
	}
	if length == 5 {
		return true
	}
	//斜着找
	length = 1
	i := x - 1
	j := y - 1
	for i >= 0 && j >= 0 {
		if a[i][j] == color {
			length++
		} else {
			break
		}
		i--
		j--
	}
	i = x + 1
	j = y + 1
	for i <= 14 && j <= 14 {
		if a[i][j] == color {
			length++
		} else {
			break
		}
		i++
		j++
	}
	if length == 5 {
		return true
	}
	return false
}
func CheckWin(a [][]int, color int, x int, y int) bool {
	var length int
	//横着找
	length = 1
	for i := x - 1; i >= 0; i-- {
		if a[i][y] == color {
			length++
		} else {
			break
		}
	}
	for i := x + 1; i <= 14; i++ {
		if a[i][y] == color {
			length++
		} else {
			break
		}
	}
	if length >= 5 {
		return true
	}
	//竖着找
	length = 1
	for i := y - 1; i >= 0; i-- {
		if a[x][i] == color {
			length++
		} else {
			break
		}
	}
	for i := y + 1; i <= 14; i++ {
		if a[x][i] == color {
			length++
		} else {
			break
		}
	}
	if length >= 5 {
		return true
	}
	//斜着找
	length = 1
	i := x - 1
	j := y - 1
	for i >= 0 && j >= 0 {
		if a[i][j] == color {
			length++
		} else {
			break
		}
		i--
		j--
	}
	i = x + 1
	j = y + 1
	for i <= 14 && j <= 14 {
		if a[i][j] == color {
			length++
		} else {
			break
		}
		i++
		j++
	}
	if length >= 5 {
		return true
	}
	return false
}
