package router

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"gobang/app/api/internal/dao/jwt"
	"gobang/app/api/internal/service"
)

func InitRouter(port string) error {
	r := gin.Default()
	r.Use(static.Serve("/", static.LocalFile("./dist(6)", false)))

	r.POST("/registerPackage", service.Register)
	r.POST("/loginPackage", service.Login)
	r.GET("/checkLoginPackage", jwt.Test)
	r.POST("/getInformationPackage", service.GetUserInformation)
	r.POST("/editInformationPackage", service.EditUserInformation)
	r.POST("/getRoom", service.GetRoom)
	r.POST("/gameHall/creatRoom", service.CreatRoom)
	v1 := r.Group("").Use() //jwt.CheckToken()
	//v1.POST("/gameHall/creatRoom", service.CreatRoom)
	v1.POST("/gameHall/enterRoom", service.EnterRoom)
	r.NoRoute(func(c *gin.Context) {
		// 在这里处理没有匹配到路由的情况，可以返回默认的文件或处理程序
		c.File("./dist(6)/index.html")
	})
	err := r.Run(":" + port)
	if err != nil {
		return err
	}
	return err
}
