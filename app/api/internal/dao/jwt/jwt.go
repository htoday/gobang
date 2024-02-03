package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gobang/app/api/global"
	"gobang/app/api/internal/model"
	"time"
)

func UseJWT(username string) string {
	MySigningKey := []byte("666")

	c := model.MyClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 60,
			ExpiresAt: time.Now().Unix() + 60*60*6,
			Issuer:    "htoday",
		},
	}
	MyToken := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	str, err := MyToken.SignedString(MySigningKey)
	if err != nil {
		global.Logger.Warn("token signed failed")
		return str
	}
	return str
}
func CheckToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		MySigningKey := []byte("666")
		str := c.GetHeader("token")
		token, err := jwt.ParseWithClaims(str, &model.MyClaims{}, func(token *jwt.Token) (interface{}, error) {
			return MySigningKey, nil
		})
		if err != nil {
			global.Logger.Warn(token.Claims.(*model.MyClaims).Username + "'s token wrong or expired" + err.Error())
			c.JSON(401, gin.H{
				"msg": "token wrong or expired",
			})
			return
		}
		global.Logger.Info(token.Claims.(*model.MyClaims).Username + "'s token correct")
		c.Next()
	}
}
func Test(c *gin.Context) {

	c.JSON(200, gin.H{
		"msg": "test",
	})
}
