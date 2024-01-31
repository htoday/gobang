package router

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"gobang/app/api/internal/dao/jwt"
	"gobang/app/api/internal/service"
)

func InitRouter(port string) error {
	r := gin.Default()
	r.Use(static.Serve("/", static.LocalFile("./dist(30)", false)))
	r.POST("/registerPackage", service.Register)
	r.POST("/loginPackage", service.Login)
	r.GET("/checkLoginPackage", jwt.Test)
	r.POST("/getInformationPackage", service.GetUserInformation)
	r.POST("/editInformationPackage", service.EditUserInformation)
	r.POST("/getRoomListPackage", service.GetRoomList)
	r.POST("/getRoomInformationPackage", service.GetRoomInformation)
	r.POST("/createRoomPackage", service.CreatRoom)
	r.POST("/enterRoomPackage", service.EnterRoom)
	r.POST("/leaveRoomPackage", service.LeaveRoom)
	r.POST("/putChessPackage", service.PutChess)
	r.POST("/startGamePackage", service.StartGame)
	r.POST("/changeFirstActionPackage", service.SetFirstAct)
	r.POST("/playerReadyPackage", service.RoomReady)
	r.POST("/changeForbiddenPackage", service.ChangeForbidden)
	//v1 := r.Group("").Use() //jwt.CheckToken()
	//v1.POST("/gameHall/creatRoom", service.CreatRoom)

	r.NoRoute(func(c *gin.Context) {
		// 在这里处理没有匹配到路由的情况，可以返回默认的文件或处理程序
		c.File("./dist(30)/index.html")
	})
	err := r.Run(":" + port)
	if err != nil {
		return err
	}
	return err
}
