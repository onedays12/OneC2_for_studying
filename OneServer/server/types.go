package server

import (
	"OneServer/controller"
	"OneServer/handler"
	"OneServer/profile"
	"OneServer/utils/response"
	"OneServer/utils/safeType"
)

type TeamServer struct {
	Profile    *profile.Profile
	Controller *controller.Controller
	Listeners  safeType.Map
	Beacons    safeType.Map
	Handler    *handler.Handler
}

type Beacon struct {
	Data       response.BeaconData
	Tick       bool
	Active     bool
	TasksQueue *safeType.Slice
}

const (
	TYPE_TASK = 1
)
