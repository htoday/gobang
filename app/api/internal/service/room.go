package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"gobang/app/api/global"
	"gobang/app/api/internal/dao/mysql"
	"gobang/app/api/internal/model"
	"strconv"
	"time"
)

type Room struct {
	RoomID        int    `json:"roomID"`
	Title         string `json:"title"`
	User1         string `json:"user1"`
	User2         string `json:"user2"`
	Forbidden     bool   `json:"forbidden"`
	Ranking       bool   `json:"ranking"`
	RoomPassword  string `json:"roomPassword"`
	User1Ready    bool   `json:"user1Ready"`
	User2Ready    bool   `json:"user2Ready"`
	RoomStatus    int    `json:"roomStatus"`
	UserNickname1 string `json:"userNickname1"`
	UserNickname2 string `json:"userNickname2"`
	FirstAct      int    `json:"firstAct"` //1为玩家1先手，2为玩家2先手，3为随机
}

func GetRoom(roomID int) (Room, error) {
	var thisRoom Room
	ctx := context.Background()
	roomJSON, err := global.RDB.HGet(ctx, "rooms", strconv.Itoa(roomID)).Result()
	if err != nil {
		global.Logger.Warn("translate JSON into room failed," + err.Error())
		return Room{}, err
	}
	//反序列化
	err = json.Unmarshal([]byte(roomJSON), &thisRoom)
	if err != nil {
		global.Logger.Warn("Unmarshal JSON failed," + err.Error())
		return thisRoom, err
	}
	return thisRoom, nil
}
func UpdateRoom(r Room, roomID int) error {
	ctx := context.Background()
	roomJSON, err := json.Marshal(r)
	if err != nil {
		return err
	}
	err = global.RDB.HSet(ctx, "rooms", strconv.Itoa(roomID), roomJSON).Err()
	if err != nil {
		return err
	}
	return nil
}
func GetRoomList(c *gin.Context) {
	var m model.Message
	err := c.ShouldBindJSON(&m)
	if err != nil {
		global.Logger.Warn("get request failed," + err.Error())
		c.JSON(400, gin.H{
			"msg": "get request failed," + err.Error(),
		})
		return
	}
	ctx := context.Background()
	roomsMap, err := global.RDB.HGetAll(ctx, "rooms").Result()
	if err != nil {
		global.Logger.Warn("get room information failed," + err.Error())
		c.JSON(400, gin.H{
			"msg": "get room information failed," + err.Error(),
		})
		return
	}
	var rooms []Room
	joinedRoomID := 0
	for _, roomJSON := range roomsMap {
		var thisRoom Room
		err := json.Unmarshal([]byte(roomJSON), &thisRoom)
		if err != nil {
			global.Logger.Warn("translate JSON into room failed," + err.Error())
			//这里不return了
		}
		thisRoom.RoomPassword = "unknown"
		if thisRoom.User1 == m.Username || thisRoom.User2 == m.Username {
			joinedRoomID = thisRoom.RoomID
		}
		rooms = append(rooms, thisRoom)
	}
	global.Logger.Info("get room information success")
	c.JSON(200, gin.H{
		"joinedRoomID": joinedRoomID,
		"room":         rooms,
	})
	//这里后面需要讨论
}
func CreatRoom(c *gin.Context) {
	var m model.Message
	err := c.ShouldBindJSON(&m)
	if err != nil {
		global.Logger.Warn("get request failed," + err.Error())
		c.JSON(400, gin.H{
			"msg": "get request failed," + err.Error(),
		})
		return
	}
	//判断名字为空的情况
	if m.Username == "" {
		global.Logger.Warn("username cannot be empty")
		c.JSON(400, gin.H{
			"msg": "username cannot empty",
		})
		return
	}
	nickname, err := mysql.GetNickname(m.Username)
	if err != nil {
		global.Logger.Warn("get nickname failed")
		c.JSON(400, gin.H{
			"msg": "get nickname failed",
		})
		return
	}
	//分配房间id
	global.Counter.Mu.Lock()
	global.Counter.Value++
	roomID := global.Counter.Value
	global.Counter.Mu.Unlock()
	r := Room{
		RoomID:        roomID,
		User1:         m.Username,
		User2:         "",
		UserNickname1: nickname,
		UserNickname2: "",
		Ranking:       m.Ranking,
		RoomPassword:  m.RoomPassword,
		Forbidden:     m.Forbidden,
		RoomStatus:    0,
		User1Ready:    false,
		User2Ready:    false,
		Title:         m.Title,
		FirstAct:      1,
	}
	err = UpdateRoom(r, r.RoomID)
	if err != nil {
		global.Logger.Warn("save room information failed," + err.Error())
		c.JSON(400, gin.H{
			"msg": "save room information failed," + err.Error(),
		})
		return
	}
	global.Logger.Info("save room information success,")
	c.JSON(200, gin.H{
		"roomID": r.RoomID,
		"msg":    "save room information success,",
	})
}
func EnterRoom(c *gin.Context) {
	var m model.Message
	err := c.ShouldBindJSON(&m)
	if err != nil {
		global.Logger.Warn("get roomID failed," + err.Error())
		c.JSON(400, gin.H{
			"msg": "get roomID failed," + err.Error(),
		})
		return
	}
	//判断名字为空的情况
	if m.Username == "" {
		global.Logger.Warn("username cannot empty")
		c.JSON(400, gin.H{
			"msg": "username cannot empty",
		})
		return
	}
	ctx := context.Background()
	roomID := m.RoomID
	//序列化
	roomJSON, err := global.RDB.HGet(ctx, "rooms", strconv.Itoa(roomID)).Result()
	if err != nil {
		global.Logger.Warn("translate JSON into room failed," + err.Error())
		c.JSON(400, gin.H{
			"msg": "translate JSON into room failed," + err.Error(),
		})
		return
	}
	var thisRoom Room
	//反序列化
	err = json.Unmarshal([]byte(roomJSON), &thisRoom)
	if err != nil {
		global.Logger.Warn("Unmarshal JSON failed," + err.Error())
		c.JSON(400, gin.H{
			"msg": "Unmarshal JSON failed," + err.Error(),
		})
		return
	}
	if thisRoom.RoomPassword != m.RoomPassword {
		global.Logger.Warn("password wrong")
		c.JSON(400, gin.H{
			"msg": "password wrong",
		})
		return
	}
	//用户名相同
	if thisRoom.User1 == m.Username || thisRoom.User2 == m.Username {
		global.Logger.Warn(m.Username + " have entered room")
		c.JSON(404, gin.H{
			"msg": m.Username + " have entered room",
		})
		return
	}
	//房间里没有人的情况
	if thisRoom.User1 == "" && thisRoom.User2 == "" {
		global.Logger.Warn("Room not exist")
		c.JSON(404, gin.H{
			"msg": "Room not exist",
		})
		return
	}
	if thisRoom.User1 != "" && thisRoom.User2 == "" {
		thisRoom.User2 = m.Username
		thisRoom.UserNickname2, err = mysql.GetNickname(thisRoom.User2)
		if err != nil {
			global.Logger.Warn("get nickname failed" + err.Error())
			c.JSON(400, gin.H{
				"msg": "get nickname failed" + err.Error(),
			})
			return
		}
		err = UpdateRoom(thisRoom, thisRoom.RoomID)
		if err != nil {
			global.Logger.Warn("save room information failed," + err.Error())
			c.JSON(400, gin.H{
				"msg": "save room information failed," + err.Error(),
			})
			return
		}
		global.Logger.Info("save room information success")
		c.JSON(200, gin.H{
			"msg": "save room information success",
		})
		return
	}
	if thisRoom.User1 != "" && thisRoom.User2 != "" {
		global.Logger.Warn("room is full")
		c.JSON(400, gin.H{
			"msg": "room is full",
		})
		return
	}
}
func LeaveRoom(c *gin.Context) {
	var m model.Message
	ctx := context.Background()
	err := c.ShouldBindJSON(&m)
	if err != nil {
		global.Logger.Warn("get roomID failed," + err.Error())
		c.JSON(400, gin.H{
			"msg": "get roomID failed," + err.Error(),
		})
		return
	}
	//判断名字为空的情况
	if m.Username == "" {
		global.Logger.Warn("username cannot empty")
		c.JSON(400, gin.H{
			"msg": "username cannot empty",
		})
		return
	}
	roomID := m.RoomID
	roomJSON, err := global.RDB.HGet(ctx, "rooms", strconv.Itoa(roomID)).Result()
	if err != nil {
		global.Logger.Warn("translate JSON into room failed," + err.Error())
		c.JSON(400, gin.H{
			"msg": "translate JSON into room failed," + err.Error(),
		})
		return
	}
	//反序列化
	var thisRoom Room
	err = json.Unmarshal([]byte(roomJSON), &thisRoom)
	if err != nil {
		global.Logger.Warn("Unmarshal JSON failed," + err.Error())
		c.JSON(400, gin.H{
			"msg": "Unmarshal JSON failed," + err.Error(),
		})
		return
	}
	if thisRoom.User1 != m.Username && thisRoom.User2 != m.Username {
		global.Logger.Warn(m.Username + "not in this room")
		c.JSON(400, gin.H{
			"msg": m.Username + "not in this room",
		})
		return
	}
	if thisRoom.User1 == m.Username {
		thisRoom.User1 = ""
		thisRoom.UserNickname1 = ""
		err = UpdateRoom(thisRoom, thisRoom.RoomID)
	}
	if thisRoom.User2 == m.Username {
		thisRoom.User2 = ""
		thisRoom.UserNickname2 = ""
		if thisRoom.RoomStatus == 1 && thisRoom.Ranking == true {
			go ResetRoom(thisRoom.RoomID)
			mysql.UpdateStar(thisRoom.User1, thisRoom.User2)
		}
		err = UpdateRoom(thisRoom, thisRoom.RoomID)
	}
	//房主为空就删房间
	if thisRoom.User1 == "" {
		err = DelRoom(thisRoom)
		if err != nil {
			global.Logger.Warn("Delete room failed," + err.Error())
			c.JSON(400, gin.H{
				"msg": "Delete room failed," + err.Error(),
			})
			return
		}
	}
	global.Logger.Info(m.Username + "leave " + strconv.Itoa(roomID) + " success")
	c.JSON(200, gin.H{
		"msg": m.Username + "leave room success",
	})
	return
}
func DelRoom(thisRoom Room) error {
	ctx := context.Background()
	err := global.RDB.HDel(ctx, "rooms", strconv.Itoa(thisRoom.RoomID)).Err()
	if err != nil {
		return err
	}
	err = global.RDB.HDel(ctx, "games", strconv.Itoa(thisRoom.RoomID)).Err()
	if err != nil {
		return err
	}
	return nil
}
func RoomReady(c *gin.Context) {
	var m model.Message
	ctx := context.Background()
	err := c.ShouldBindJSON(&m)
	if err != nil {
		global.Logger.Warn("get roomID failed," + err.Error())
		c.JSON(400, gin.H{
			"msg": "get roomID failed," + err.Error(),
		})
		return
	}
	//判断名字为空的情况
	if m.Username == "" {
		global.Logger.Warn("username cannot empty")
		c.JSON(400, gin.H{
			"msg": "username cannot empty",
		})
		return
	}
	roomID := m.RoomID
	roomJSON, err := global.RDB.HGet(ctx, "rooms", strconv.Itoa(roomID)).Result()
	if err != nil {
		global.Logger.Warn("translate JSON into room failed," + err.Error())
		c.JSON(400, gin.H{
			"msg": "translate JSON into room failed," + err.Error(),
		})
		return
	}
	//反序列化
	var thisRoom Room
	err = json.Unmarshal([]byte(roomJSON), &thisRoom)
	if err != nil {
		global.Logger.Warn("Unmarshal JSON failed," + err.Error())
		c.JSON(400, gin.H{
			"msg": "Unmarshal JSON failed," + err.Error(),
		})
		return
	}
	if thisRoom.User2 == m.Username {
		if thisRoom.User2Ready == false {
			thisRoom.User2Ready = true
			err = UpdateRoom(thisRoom, thisRoom.RoomID)
			if err != nil {
				global.Logger.Warn("Update Room Ready failed," + err.Error())
				c.JSON(400, gin.H{
					"msg": "Update Room Ready failed," + err.Error(),
				})
				return
			}
			global.Logger.Info(m.Username + "ready success")
			c.JSON(200, gin.H{
				"msg": m.Username + "ready success",
			})
			return
		} else {
			thisRoom.User2Ready = false
			err = UpdateRoom(thisRoom, thisRoom.RoomID)
			if err != nil {
				global.Logger.Warn("Update Room Ready failed," + err.Error())
				c.JSON(400, gin.H{
					"msg": "Update Room Ready failed," + err.Error(),
				})
				return
			}
			global.Logger.Info(m.Username + "ready success")
			c.JSON(200, gin.H{
				"msg": m.Username + "ready success",
			})
			return
		}

	}

	if thisRoom.User1 == m.Username {
		if thisRoom.User1Ready == false {
			thisRoom.User1Ready = true
			err = UpdateRoom(thisRoom, thisRoom.RoomID)
			if err != nil {
				global.Logger.Warn("Update Room Ready failed," + err.Error())
				c.JSON(400, gin.H{
					"msg": "Update Room Ready failed," + err.Error(),
				})
				return
			}
			global.Logger.Info(m.Username + "ready success")
			c.JSON(200, gin.H{
				"msg": m.Username + "ready success",
			})
			return
		}
		if thisRoom.User1Ready == true {
			thisRoom.User1Ready = false
			err = UpdateRoom(thisRoom, thisRoom.RoomID)
			if err != nil {
				global.Logger.Warn("Update Room Ready failed," + err.Error())
				c.JSON(400, gin.H{
					"msg": "Update Room Ready failed," + err.Error(),
				})
				return
			}
			global.Logger.Info(m.Username + " cansel ready success")
			c.JSON(200, gin.H{
				"msg": m.Username + " cansel ready success",
			})
			return
		}
	}
	global.Logger.Warn("username wrong")
	c.JSON(400, gin.H{
		"msg": "username wrong",
	})
}

func ChangeForbidden(c *gin.Context) {
	var m model.Message
	ctx := context.Background()
	err := c.ShouldBindJSON(&m)
	if err != nil {
		global.Logger.Warn("get roomID failed," + err.Error())
		c.JSON(400, gin.H{
			"msg": "get roomID failed," + err.Error(),
		})
		return
	}
	//判断名字为空的情况
	if m.Username == "" {
		global.Logger.Warn("username cannot empty")
		c.JSON(400, gin.H{
			"msg": "username cannot empty",
		})
		return
	}
	roomID := m.RoomID
	roomJSON, err := global.RDB.HGet(ctx, "rooms", strconv.Itoa(roomID)).Result()
	if err != nil {
		global.Logger.Warn("translate JSON into room failed," + err.Error())
		c.JSON(400, gin.H{
			"msg": "translate JSON into room failed," + err.Error(),
		})
		return
	}
	//反序列化
	var thisRoom Room
	err = json.Unmarshal([]byte(roomJSON), &thisRoom)
	if err != nil {
		global.Logger.Warn("Unmarshal JSON failed," + err.Error())
		c.JSON(400, gin.H{
			"msg": "Unmarshal JSON failed," + err.Error(),
		})
		return
	}
	//这里商讨一下该不该确认房主
	if thisRoom.Forbidden == true {
		thisRoom.Forbidden = false

	} else {
		thisRoom.Forbidden = true
	}
	err = UpdateRoom(thisRoom, thisRoom.RoomID)
	if err != nil {
		global.Logger.Warn("Update Room Ready failed," + err.Error())
		c.JSON(400, gin.H{
			"msg": "Update Room Ready failed," + err.Error(),
		})
		return
	}
	global.Logger.Info("forbidden changed in room " + strconv.Itoa(thisRoom.RoomID))
	c.JSON(200, gin.H{
		"msg": "forbidden changed in room " + strconv.Itoa(thisRoom.RoomID),
	})
}
func SetFirstAct(c *gin.Context) {
	var m model.Message
	err := c.ShouldBindJSON(&m)
	if err != nil {
		global.Logger.Warn("get roomID failed," + err.Error())
		c.JSON(404, gin.H{
			"msg": "get roomID failed," + err.Error(),
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
	if thisRoom.FirstAct == 1 {
		thisRoom.FirstAct = 2
	} else if thisRoom.FirstAct == 2 {
		thisRoom.FirstAct = 3
	} else if thisRoom.FirstAct == 3 {
		thisRoom.FirstAct = 1
	}
	err = UpdateRoom(thisRoom, thisRoom.RoomID)
	if err != nil {
		global.Logger.Warn("Update Room Ready failed," + err.Error())
		c.JSON(400, gin.H{
			"msg": "Update Room Ready failed," + err.Error(),
		})
		return
	}
	global.Logger.Info("Set FirstAct success")
	c.JSON(200, gin.H{
		"msg": "Set FirstAct success",
	})
	return
}
func GetRoomInformation(c *gin.Context) {
	var m model.Message
	ctx := context.Background()
	err := c.ShouldBindJSON(&m)
	if err != nil {
		global.Logger.Warn("get roomID failed," + err.Error())
		c.JSON(404, gin.H{
			"msg": "get roomID failed," + err.Error(),
		})
		return
	}
	roomID := m.RoomID
	roomJSON, err := global.RDB.HGet(ctx, "rooms", strconv.Itoa(roomID)).Result()
	if err != nil {
		global.Logger.Warn("translate JSON into room failed," + err.Error())
		c.JSON(404, gin.H{
			"msg": "translate JSON into room failed," + err.Error(),
		})
		return
	}
	//反序列化
	var thisRoom Room
	err = json.Unmarshal([]byte(roomJSON), &thisRoom)
	if err != nil {
		global.Logger.Warn("Unmarshal JSON failed," + err.Error())
		c.JSON(400, gin.H{
			"msg": "Unmarshal JSON failed," + err.Error(),
		})
		return
	}
	if thisRoom.RoomStatus == 0 {
		c.JSON(200, gin.H{
			"roomID":       thisRoom.RoomID,
			"username1":    thisRoom.User1,
			"nickname1":    thisRoom.UserNickname1,
			"username2":    thisRoom.User2,
			"nickname2":    thisRoom.UserNickname2,
			"forbidden":    thisRoom.Forbidden,
			"roomStatus":   thisRoom.RoomStatus,
			"roomPassword": thisRoom.RoomPassword,
			"ranking":      thisRoom.Ranking,
			"user1Ready":   thisRoom.User1Ready,
			"user2Ready":   thisRoom.User2Ready,
			"firstAct":     thisRoom.FirstAct,
			"title":        thisRoom.Title,
			"msg":          "get room information success",
		})
		return
	}
	thisGame, err := GetGame(thisRoom.RoomID)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "get game JSON failed," + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"roomID":       thisRoom.RoomID,
		"username1":    thisRoom.User1,
		"nickname1":    thisRoom.UserNickname1,
		"username2":    thisRoom.User2,
		"nickname2":    thisRoom.UserNickname2,
		"forbidden":    thisRoom.Forbidden,
		"roomStatus":   thisRoom.RoomStatus,
		"roomPassword": thisRoom.RoomPassword,
		"ranking":      thisRoom.Ranking,
		"user1Ready":   thisRoom.User1Ready,
		"user2Ready":   thisRoom.User2Ready,
		"firstAct":     thisRoom.FirstAct,
		"gameBoard":    thisGame.CheckBoard,
		"turn":         thisGame.Turn,
		"msg":          "get room information success",
	})
}
func ResetRoom(roomID int) {
	time.Sleep(time.Second * 4)
	thisRoom, err := GetRoom(roomID)
	if err != nil {
		global.Logger.Error("get room failed" + err.Error())
	}
	thisRoom.RoomStatus = 0
	thisGame, err := GetGame(roomID)
	if err != nil {
		global.Logger.Error("get game failed" + err.Error())
	}
	err = UpdateRoom(thisRoom, thisRoom.RoomID)
	err = UpdateGame(thisGame, thisGame.RoomID)
}
func ChangeRanking(c *gin.Context) {
	var m model.Message
	ctx := context.Background()
	err := c.ShouldBindJSON(&m)
	if err != nil {
		global.Logger.Warn("get roomID failed," + err.Error())
		c.JSON(400, gin.H{
			"msg": "get roomID failed," + err.Error(),
		})
		return
	}
	//判断名字为空的情况
	if m.Username == "" {
		global.Logger.Warn("username cannot empty")
		c.JSON(400, gin.H{
			"msg": "username cannot empty",
		})
		return
	}
	roomID := m.RoomID
	roomJSON, err := global.RDB.HGet(ctx, "rooms", strconv.Itoa(roomID)).Result()
	if err != nil {
		global.Logger.Warn("translate JSON into room failed," + err.Error())
		c.JSON(400, gin.H{
			"msg": "translate JSON into room failed," + err.Error(),
		})
		return
	}
	//反序列化
	var thisRoom Room
	err = json.Unmarshal([]byte(roomJSON), &thisRoom)
	if err != nil {
		global.Logger.Warn("Unmarshal JSON failed," + err.Error())
		c.JSON(400, gin.H{
			"msg": "Unmarshal JSON failed," + err.Error(),
		})
		return
	}

	if thisRoom.Ranking == true {
		thisRoom.Ranking = false

	} else {
		thisRoom.Ranking = true
	}
	err = UpdateRoom(thisRoom, thisRoom.RoomID)
	if err != nil {
		global.Logger.Warn("Update Room Ready failed," + err.Error())
		c.JSON(400, gin.H{
			"msg": "Update Room Ready failed," + err.Error(),
		})
		return
	}
	global.Logger.Info("ranking changed in room " + strconv.Itoa(thisRoom.RoomID))
	c.JSON(200, gin.H{
		"msg": "ranking changed in room " + strconv.Itoa(thisRoom.RoomID),
	})
}
func RoomChat() {

}
