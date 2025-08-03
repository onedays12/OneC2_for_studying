package server

import (
	"OneServer/utils/request"
	"OneServer/utils/response"
	"errors"
	"fmt"
)

func (ts *TeamServer) ListenerStart(listenerName string, configType string, config request.ConfigDetail) error {

	if ts.Listeners.Contains(listenerName) {
		return errors.New("listener already exists")
	}

	listenerData, err := ts.Handler.ListenerStart(listenerName, configType, config)
	if err != nil {
		return err
	}

	listenerData.Name = listenerName
	listenerData.Type = configType

	ts.Listeners.Put(listenerName, listenerData)

	return nil
}

func (ts *TeamServer) ListenerGetConfig(listenerName string, configType string) (request.ConfigDetail, error) {
	if !ts.Listeners.Contains(listenerName) {
		return request.ConfigDetail{}, fmt.Errorf("listener %v does not exist", listenerName)
	}

	value, _ := ts.Listeners.Get(listenerName)

	return value.(response.ListenerData).Data, nil
}
