package server

import (
	"OneServer/controller"
	"OneServer/handler"
	"OneServer/logs"
	"OneServer/profile"
	"OneServer/utils/safeType"
	"encoding/json"
	"go.uber.org/zap"
	"os"
)

func NewTeamServer() *TeamServer {
	ts := &TeamServer{
		Profile:   profile.NewProfile(),
		Listeners: safeType.NewMap(),
		Beacons:   safeType.NewMap(),
	}

	ts.Handler = handler.NewHandler(ts)
	return ts

}

func (ts *TeamServer) SetProfile(path string) error {

	fileContent, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(fileContent, &ts.Profile)
	if err != nil {
		return err
	}

	return nil
}

func (ts *TeamServer) Start() {

	var stopped = make(chan bool) // 初始化 stopped 通道

	ts.Controller = controller.NewController(ts, *ts.Profile)

	go ts.Controller.StartServer(&stopped)

	logs.Logger.Info("OneServer started -> ", zap.String("address",
		ts.Profile.TeamServerConfig.Host),
		zap.Int("port", ts.Profile.TeamServerConfig.Port))

	// 等待服务启动或超时
	<-stopped
	logs.Logger.Info("Server stopped gracefully")

	// 关闭服务器
	os.Exit(0)
}
