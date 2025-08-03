package listener

import (
	"OneServer/utils/request"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TeamServer interface {
	BeaconIsExists(beaconId string) bool
	BeaconCreate(beaconId string, beat []byte, listenerName string, ExternalIP string, Async bool) error
	BeaconGetAllTasks(beaconId string, maxDataSize int) ([]byte, error)
	BeaconProcessData(beaconId string, bodyData []byte) error
}

type ListenerHTTP struct {
	ts TeamServer
}

type HTTP struct {
	GinEngine *gin.Engine
	Server    *http.Server
	Config    request.ConfigDetail
	Name      string
	Active    bool
}
