package router

import (
	"github.com/gin-gonic/gin"
	"gobang/app/api/internal/service"
)

func InitRouter(port string) error {
	r := gin.Default()
	r.POST("/register", service.Register)
	r.POST("/login", service.Login)
	r.POST("/gameHall/creatRoom", service.CreatRoom)
	v1 := r.Group("").Use() //jwt.CheckToken()
	//v1.POST("/gameHall/creatRoom", service.CreatRoom)
	v1.POST("/gameHall/enterRoom", service.EnterRoom)

	err := r.Run(":" + port)
	if err != nil {
		return err
	}
	return err
}
