package handler

import (
	"OneServer/utils/request"
	"OneServer/utils/response"
	"errors"
)

func (h *Handler) ListenerStart(listenerName string, configType string, config request.ConfigDetail) (response.ListenerData, error) {
	listenerHandler, ok := h.ListenerHandlers[configType]
	if !ok {
		return response.ListenerData{}, errors.New("handler not found")
	}

	err := listenerHandler.ListenerValid(config)
	if err != nil {
		return response.ListenerData{}, err
	}

	return listenerHandler.ListenerStart(listenerName, config)
}
