package routes

import (
	"github.com/gin-gonic/gin"
	"controllers"
)

func InitRoute(router  *gin.Engine){
	// 拉取用户信息
	router.GET("/getInfo", controllers.CheckUserLoginInfo())
}
