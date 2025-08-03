package controller

import (
	"OneServer/logs"
	"OneServer/utils/request"
	"OneServer/utils/response"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func (c *Controller) BeaconGenerate(ctx *gin.Context) {

	var (
		beaconConfig   request.BeaconConfig
		listenerConfig request.ConfigDetail
		err            error
	)

	err = ctx.ShouldBindJSON(&beaconConfig)
	if err != nil {
		logs.Logger.Error("Error in binding JSON data: ", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"code": false, "message": err.Error()})
		return
	}

	listenerConfig, err = c.TeamServer.ListenerGetConfig(beaconConfig.ListenerName, beaconConfig.ListenerType)
	if err != nil {
		logs.Logger.Error("Error in getting listener config: ", zap.Error(err))
		ctx.JSON(http.StatusOK, gin.H{"code": false, "message": err.Error()})
		return
	}

	fileContent, fileName, err := c.TeamServer.BeaconGenerate(beaconConfig, listenerConfig)
	if err != nil {
		logs.Logger.Error("Error in generating beacon: ", zap.Error(err))
		ctx.JSON(http.StatusOK, gin.H{"code": false, "message": err.Error()})
		return
	}

	base64FileName := base64.StdEncoding.EncodeToString([]byte(fileName))
	base64Content := base64.StdEncoding.EncodeToString(fileContent)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    true,
		"message": "Beacon generated successfully",
		"data": response.BeaconFile{
			FileName:    base64FileName,
			FileContent: base64Content,
		},
	})
}

func (c *Controller) BeaconCommandExecute(ctx *gin.Context) {
	var (
		username    string
		commandData request.CommandData
		err         error
	)

	err = ctx.ShouldBindJSON(&commandData)
	if err != nil {
		logs.Logger.Error("Error in binding JSON data: ", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"code": false, "message": err.Error()})
		return
	}

	username = "oneday"

	err = c.TeamServer.BeaconCommand(commandData.BeaconType, commandData.BeaconId, username, commandData.CmdLine, commandData.Data)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "ok": false})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Beacon command task submitted successfully", "ok": true})
}
