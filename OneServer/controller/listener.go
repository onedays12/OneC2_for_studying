package controller

import (
	"OneServer/logs"
	"OneServer/utils/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"regexp"
)

func (c *Controller) ListenerStart(ctx *gin.Context) {

	var (
		listenerConfig request.ListenerConfig
		err            error
	)

	err = ctx.ShouldBindJSON(&listenerConfig)
	if err != nil {
		logs.Logger.Error("Error in binding JSON data: ", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"code": false, "message": err.Error()})
		return
	}

	if ValidListenerName(listenerConfig.ListenerName) == false {
		logs.Logger.Error("Invalid listener name", zap.String("listener_name", listenerConfig.ListenerName))
		ctx.JSON(http.StatusOK, gin.H{"code": false, "message": "Invalid listener name"})
		return
	}

	err = c.TeamServer.ListenerStart(listenerConfig.ListenerName, listenerConfig.ConfigType, listenerConfig.Config)
	if err != nil {
		logs.Logger.Error("Error in starting listener", zap.Error(err))
		ctx.JSON(http.StatusOK, gin.H{"code": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Listener started successfully", "ok": true})
}

// 辅助函数

func ValidListenerName(s string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9-_]+$")
	return re.MatchString(s)
}
