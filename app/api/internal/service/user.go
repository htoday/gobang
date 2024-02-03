package service

import (
	"github.com/gin-gonic/gin"
	"gobang/app/api/global"
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
		c.JSON(400, gin.H{
			//"code": consts.ShouldBindFailed,
			"msg": "user register bind parameter failed," + err.Error(),
		})
		return
	}
	//将username加入redis
	flag, err := redis.CheckUser(u.Username)
	if err != nil {
		c.JSON(400, gin.H{
			//"code": consts.CheckUserFailed,
			"msg": "CheckUser Failed in redis," + err.Error(),
		})
		return
	}
	if flag == true { //如果重复就return
		c.JSON(409, gin.H{
			//"code": consts.UserAlreadyExists,
			"msg": "User Already Exists",
		})
		return
	}
	//再去mysql里面查重（根本没有必要）
	flag, err = mysql.CheckUser(u.Username)
	if err != nil {
		c.JSON(400, gin.H{
			//"code": consts.CheckUserFailed,
			"msg": "CheckUser Failed in mysql" + err.Error(),
		})
		return
	}
	if flag == true {
		c.JSON(409, gin.H{
			//"code": consts.UserAlreadyExists,
			"msg": "User Already Exists",
		})
		global.Logger.Warn("username" + u.Username + " existed in mysql")
		return
	}
	//将用户加入redis
	err = redis.AddUser(u.Username, u.Password)
	if err != nil {
		c.JSON(400, gin.H{
			//"code": consts.InsertUserFailed,
			"msg": "Insert User into redis Failed",
		})
		return
	}
	//将用户加入mysql
	err, newID := mysql.AddNewUser(u.Username, u.Password, u.Nickname)
	if err != nil {
		c.JSON(400, gin.H{
			//"code": consts.InsertUserFailed,
			"msg": "insert user data into mysql failed",
		})
		global.Logger.Warn("insert user data into mysql failed" + err.Error())
		return
	}
	c.JSON(200, gin.H{
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
		c.JSON(400, gin.H{
			//"code": consts.ShouldBindFailed,
			"msg": "user register bind parameter failed," + err.Error(),
		})
		return
	}
	//首先判断用户名是否存在
	flag, err := redis.CheckUser(u.Username)
	if err != nil {
		c.JSON(400, gin.H{
			//"code": consts.CheckUserFailed,
			"msg": "CheckUser Failed in redis," + err.Error(),
		})
		return
	}
	if flag == false { //如果用户名不存在，返回
		c.JSON(404, gin.H{
			//"code": consts.UsernameNotExists,
			"msg": "Username Not Exists",
		})
		return
	}
	//先在redis里面找
	flag, err = redis.CheckPassword(u.Username, u.Password)
	if err != nil { //如果有其他错误
		global.Logger.Warn("username" + u.Username + "redis select failed" + err.Error())
		c.JSON(400, gin.H{
			//"code": consts.CheckUserFailed,
			"msg": "redis select failed," + err.Error(),
		})
		return
	}
	if flag == true { //如果redis找到了
		InitializeUserOnline(u.Username)
		c.JSON(200, gin.H{
			//"code": consts.PasswordCorrect,
			"msg": "Password Correct,",
		})
		return
	}
	//没找到就在mysql里面找
	flag, err = mysql.CheckPassword(u.Username, u.Password)
	if err != nil { //如果有其他错误
		global.Logger.Warn("username" + u.Username + " mysql select failed " + err.Error())
		c.JSON(400, gin.H{
			//"code": consts.CheckUserFailed,
			"msg": "mysql select failed," + err.Error(),
		})
		return
	}
	if flag == true { //密码正确
		str := jwt.UseJWT(u.Username)
		c.JSON(200, gin.H{
			//"code":  consts.PasswordCorrect,
			"msg":   "Password Correct,",
			"token": str,
		})
		err1 := redis.AddUser(u.Username, u.Password)
		if err1 != nil {
			global.Logger.Warn("Password Correct But insert user into redis failed" + err1.Error())
			return
		}
		global.Logger.Info(u.Username + " login success")
		InitializeUserOnline(u.Username)
		return
	}
	c.JSON(401, gin.H{
		//"code": consts.PasswordWrong,
		"msg": "Password Wrong,",
	})
}
func GetUserInformation(c *gin.Context) {
	var u model.User
	err := c.ShouldBindJSON(&u)
	if err != nil {
		global.Logger.Warn("user register bind parameter failed," + err.Error())
		c.JSON(400, gin.H{
			"msg": "user register bind parameter failed," + err.Error(),
		})
		return
	}
	username := u.Username
	u, err = mysql.GetUserInformation(username)
	if err != nil {
		global.Logger.Warn("get information failed," + err.Error())
		c.JSON(400, gin.H{
			"msg": "get information failed," + err.Error(),
		})
	}
	global.Logger.Info("get information success,")
	c.JSON(200, gin.H{
		"username":   u.Username,
		"password":   u.Password,
		"nickname":   u.Nickname,
		"email":      u.Email,
		"uid":        u.ID,
		"starAmount": u.StarAmount,
		"msg":        "get information success",
	})
}
func EditUserInformation(c *gin.Context) {
	var u model.User
	err := c.ShouldBindJSON(&u)
	if err != nil {
		global.Logger.Warn("user register bind parameter failed," + err.Error())
		c.JSON(400, gin.H{
			"msg": "user register bind parameter failed," + err.Error(),
		})
		return
	}
	var rawUser model.User
	rawUser, err = mysql.GetUserInformation(u.Username)
	if err != nil {
		global.Logger.Warn("get information failed," + err.Error())
		c.JSON(400, gin.H{
			"msg": "get information failed," + err.Error(),
		})
	}
	var (
		newNickname string
		newPassword string
		newEmail    string
	)
	if u.Nickname != "" {
		newNickname = u.Nickname
	} else {
		newNickname = rawUser.Nickname
	}
	if u.Password != "" {
		newPassword = u.Password
	} else {
		newPassword = rawUser.Password
	}
	if u.Email != "" {
		newEmail = u.Email
	} else {
		newEmail = rawUser.Email
	}
	sqlstr := "UPDATE gobangUsers SET password=?,nickname = ?, email = ? WHERE username = ?"
	_, err = global.MysqlDB.Exec(sqlstr, newPassword, newNickname, newEmail, u.Username)
	if err != nil {
		global.Logger.Warn("edit information failed," + err.Error())
		c.JSON(400, gin.H{
			"msg": "edit information failed," + err.Error(),
		})
		return
	}
	global.Logger.Info("edit information success")
	c.JSON(200, gin.H{
		"msg": "edit information success",
	})
}
