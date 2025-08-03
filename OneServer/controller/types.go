package controller

import (
	"OneServer/utils/request"
	"github.com/gin-gonic/gin"
)

type TeamServer interface {
	ListenerStart(listenerName string, configType string, config request.ConfigDetail) error
	ListenerGetConfig(listenerName string, configType string) (request.ConfigDetail, error)
	BeaconGenerate(beaconConfig request.BeaconConfig, listenerConfig request.ConfigDetail) ([]byte, string, error)
	BeaconCommand(beaconName string, beaconId string, clientName string, cmdline string, args map[string]any) error
}

type Controller struct {
	Host     string
	Port     int
	Endpoint string
	Hash     string
	Cert     string
	Key      string

	TeamServer TeamServer
	Engine     *gin.Engine
}
