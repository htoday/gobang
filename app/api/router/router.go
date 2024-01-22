package router

import (
	"github.com/gin-gonic/gin"
	"gobang/app/api/internal/dao/jwt"
	"gobang/app/api/internal/service"
)

func InitRouter(port string) error {
	r := gin.Default()
	r.POST("/register", service.Register)
	v1 := r.Group("v1").Use(jwt.CheckToken())
	r.POST("/login", service.Login)
	v1.POST("/gameHall/creatRoom", service.CreatRoom)
	v1.POST("/gameHall/joinRoom", service.JoinRoom)

	err := r.Run(":" + port)
	if err != nil {
		return err
	}
	return err
}
