package middleware

import (
	"github.com/gin-gonic/gin"
	"models"
	)

func CheckUserInfo() gin.HandlerFunc{
	return func(c *gin.Context) {
		sessionCk := models.NewSessionCookie(c)
		retCode,_ := models.CheckUserLoginInfo(c,sessionCk)
		if retCode != 0{
			c.JSON(200, gin.H{

			})
		}
		// 用户登录态校验通过
		c.Next()
	}
}
