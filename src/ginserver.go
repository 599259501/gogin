package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"routes"
	"config"
	"fmt"
)

func main(){
	apiServer,err := NewApiServer()
	if err!=nil{
		fmt.Printf("Init Api Server has err=%v,please view NewApiServer func",err)
		return
	}
	/*
		启动服务进程
	*/
	fmt.Printf("begin start server...")
	apiServer.Start()
}

type ApiServer struct{
	Config map[string]interface{}
	Server *http.Server
	Router *gin.Engine
}

func NewApiServer() (*ApiServer,error){
	router := gin.Default()
	/*
		加载路由信息，所有的路由信息统一在routes/api.go目录下配置
		业务逻辑代码可以在controllers目录下写

	*/
	routes.InitRoute(router)
	/*
		加载配置文件,所有配置文件都在config目录下配置
		配置文件的格式使用yaml格式解析
	*/
	configArr,err := config.InitConfig()

	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return &ApiServer{
		Config:configArr,
		Server:s,
		Router:router,
	},err
}

func (apiSvr *ApiServer) Start() error{
	return apiSvr.Server.ListenAndServe()
}
