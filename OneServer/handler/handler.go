package handler

import (
	"OneServer/handler/beacon"
	"OneServer/handler/listener"
)

func NewHandler(teamserver any) *Handler {
	h := &Handler{
		ListenerHandlers: make(map[string]ListenerHandler),
		BeaconHandlers:   make(map[string]BeaconHandler),
	}

	// 创建并注册 Beacon-HTTP 监听器
	h.ListenerHandlers["Beacon-HTTP"] = listener.NewBeaconHTTPListener(teamserver).(ListenerHandler)
	h.BeaconHandlers["Beacon"] = beacon.NewBeaconHandler(teamserver).(BeaconHandler)
	return h
}
