package handler

import (
	"OneServer/utils/request"
	"OneServer/utils/response"
)

type ListenerHandler interface {
	ListenerValid(listenerConfig request.ConfigDetail) error
	ListenerStart(name string, data request.ConfigDetail) (response.ListenerData, error)
}

type BeaconHandler interface {
	BeaconGenerate(beaconConfig request.GenerateConfig, listenerConfig request.ConfigDetail) ([]byte, string, error)
	BeaconCreate(beat []byte) (response.BeaconData, error)
	BeaconCommand(client string, cmdline string, beaconData response.BeaconData, args map[string]any) error
	BeaconPackData(beaconData response.BeaconData, tasks []response.TaskData) ([]byte, error)
	BeaconProcessData(beaconData response.BeaconData, packedData []byte) ([]byte, error)
}

type Handler struct {
	ListenerHandlers map[string]ListenerHandler
	BeaconHandlers   map[string]BeaconHandler
}
