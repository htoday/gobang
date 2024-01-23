package service

import (
	"github.com/gin-gonic/gin"
	"gobang/app/api/global"
	"gobang/app/api/internal/consts"
	"gobang/app/api/internal/dao/jwt"
	"gobang/app/api/internal/dao/mysql"
	"gobang/app/api/internal/dao/redis"
	"gobang/app/api/internal/model"
)

func Register(c *gin.Context) {
	var u model.User
	err := c.ShouldBindJSON(&u)
	if err != nil {
		global.Logger.Error("user register bind parameter failed," + err.Error())
		c.JSON(consts.ShouldBindFailed, gin.H{
			//"code": consts.ShouldBindFailed,
			"msg": "user register bind parameter failed," + err.Error(),
		})
		return
	}
	//将username加入redis
	flag, err := redis.CheckUser(u.Username)
	if err != nil {
		c.JSON(consts.CheckUserFailed, gin.H{
			//"code": consts.CheckUserFailed,
			"msg": "CheckUser Failed in redis," + err.Error(),
		})
		return
	}
	if flag == true { //如果重复就return
		c.JSON(consts.UserAlreadyExists, gin.H{
			//"code": consts.UserAlreadyExists,
			"msg": "User Already Exists",
		})
		return
	}
	//再去mysql里面查重（根本没有必要）
	flag, err = mysql.CheckUser(u.Username)
	if err != nil {
		c.JSON(consts.CheckUserFailed, gin.H{
			//"code": consts.CheckUserFailed,
			"msg": "CheckUser Failed in mysql" + err.Error(),
		})
		return
	}
	if flag == true {
		c.JSON(consts.UserAlreadyExists, gin.H{
			//"code": consts.UserAlreadyExists,
			"msg": "User Already Exists",
		})
		global.Logger.Warn("username" + u.Username + " existed in mysql")
		return
	}
	//将用户加入redis
	err = redis.AddUser(u.Username, u.Password)
	if err != nil {
		c.JSON(consts.InsertUserFailed, gin.H{
			//"code": consts.InsertUserFailed,
			"msg": "Insert User into redis Failed",
		})
		return
	}
	//将用户加入mysql
	err, newID := mysql.AddNewUser(u.Username, u.Password, u.Nickname)
	if err != nil {
		c.JSON(consts.InsertUserFailed, gin.H{
			//"code": consts.InsertUserFailed,
			"msg": "insert user data into mysql failed",
		})
		global.Logger.Warn("insert user data into mysql failed" + err.Error())
		return
	}
	c.JSON(consts.RegisterSuccess, gin.H{
		//"code": consts.RegisterSuccess,
		"msg": "Register Success",
		"id":  newID,
	})
	global.Logger.Info("register" + u.Username + "success")
}

func Login(c *gin.Context) {
	var u model.User
	err := c.ShouldBindJSON(&u)
	if err != nil {
		global.Logger.Error("user register bind parameter failed," + err.Error())
		c.JSON(consts.ShouldBindFailed, gin.H{
			//"code": consts.ShouldBindFailed,
			"msg": "user register bind parameter failed," + err.Error(),
		})
		return
	}
	//首先判断用户名是否存在
	flag, err := redis.CheckUser(u.Username)
	if err != nil {
		c.JSON(consts.CheckUserFailed, gin.H{
			//"code": consts.CheckUserFailed,
			"msg": "CheckUser Failed in redis," + err.Error(),
		})
		return
	}
	if flag == false { //如果用户名不存在，返回
		c.JSON(consts.UsernameNotExists, gin.H{
			//"code": consts.UsernameNotExists,
			"msg": "Username Not Exists",
		})
		return
	}
	//先在redis里面找
	flag, err = redis.CheckPassword(u.Username, u.Password)
	if err != nil { //如果有其他错误
		c.JSON(consts.CheckUserFailed, gin.H{
			//"code": consts.CheckUserFailed,
			"msg": "redis select failed," + err.Error(),
		})
		return
	}
	if flag == true { //如果redis找到了
		c.JSON(consts.PasswordCorrect, gin.H{
			//"code": consts.PasswordCorrect,
			"msg": "Password Correct," + err.Error(),
		})
		return
	}
	//没找到就在mysql里面找
	flag, err = mysql.CheckPassword(u.Username, u.Password)
	if err != nil { //如果有其他错误
		c.JSON(consts.CheckUserFailed, gin.H{
			//"code": consts.CheckUserFailed,
			"msg": "mysql select failed," + err.Error(),
		})
		return
	}
	if flag == true { //密码正确
		str := jwt.UseJWT(u.Username)
		c.JSON(consts.PasswordCorrect, gin.H{
			//"code":  consts.PasswordCorrect,
			"msg":   "Password Correct," + err.Error(),
			"token": str,
		})
		err1 := redis.AddUser(u.Username, u.Password)
		if err1 != nil {
			global.Logger.Warn("Password Correct But insert user into redis failed" + err1.Error())
			return
		}
		global.Logger.Info(u.Username + " login success")
		return
	}
	c.JSON(consts.PasswordWrong, gin.H{
		//"code": consts.PasswordWrong,
		"msg": "Password Wrong," + err.Error(),
	})
}
