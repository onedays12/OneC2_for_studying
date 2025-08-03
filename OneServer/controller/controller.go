package controller

import (
	"OneServer/logs"
	"OneServer/middlewares"
	"OneServer/profile"
	"OneServer/utils/crypt"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"os"
)

func NewController(ts TeamServer, tsProfile profile.Profile) *Controller {

	gin.SetMode(tsProfile.TeamServerConfig.Env)

	if tsProfile.ServerResponseConfig.PagePath != "" {
		fileContent, _ := os.ReadFile(tsProfile.ServerResponseConfig.PagePath)
		tsProfile.ServerResponseConfig.Page = string(fileContent)
	}

	controller := new(Controller)
	controller.Host = tsProfile.TeamServerConfig.Host
	controller.Port = tsProfile.TeamServerConfig.Port
	controller.Endpoint = tsProfile.TeamServerConfig.Endpoint
	controller.Hash = crypt.SHA256([]byte(tsProfile.TeamServerConfig.Password))
	controller.Cert = tsProfile.TeamServerConfig.Cert
	controller.Key = tsProfile.TeamServerConfig.Key
	controller.TeamServer = ts

	router := gin.Default()

	router.Use(middlewares.GinLogger(), middlewares.GinRecovery(true))
	router.Use(middlewares.Default404Middleware(tsProfile.ServerResponseConfig))

	controller.Engine = router
	controller.InitRouter()

	return controller
}

func (c *Controller) InitRouter() {

	apiGroup := c.Engine.Group(c.Endpoint)

	apiGroup.GET("/test", c.Test)
	apiGroup.POST("/listener/create", c.ListenerStart)
	apiGroup.POST("/beacon/generate", c.BeaconGenerate)
	apiGroup.POST("/beacon/command/execute", c.BeaconCommandExecute)
}

func (c *Controller) StartServer(finished *chan bool) {

	host := fmt.Sprintf("%s:%d", c.Host, c.Port)

	err := c.Engine.RunTLS(host, c.Cert, c.Key)

	// 发送错误结束
	if err != nil {
		logs.Logger.Error("Failed to start server", zap.Error(err))
	}

	// 正常结束
	*finished <- true

}
