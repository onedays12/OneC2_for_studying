package handler

import (
	"OneServer/utils/request"
	"OneServer/utils/response"
	"errors"
)

func (h *Handler) BeaconGenerate(beaconConfig request.BeaconConfig, listenerConfig request.ConfigDetail) ([]byte, string, error) {
	beaconHandler, ok := h.BeaconHandlers[beaconConfig.BeaconType]
	if !ok {
		return nil, "", errors.New("beacon handler not found")
	}

	return beaconHandler.BeaconGenerate(beaconConfig.Config, listenerConfig)
}

func (h *Handler) BeaconCreate(beaconName string, beat []byte) (response.BeaconData, error) {
	beaconHandler, ok := h.BeaconHandlers[beaconName]
	if !ok {
		return response.BeaconData{}, errors.New("beacon handler not found")
	}

	return beaconHandler.BeaconCreate(beat)
}

func (h *Handler) BeaconCommand(client string, cmdline string, beaconName string, beaconData response.BeaconData, args map[string]any) error {
	beaconHandler, ok := h.BeaconHandlers[beaconName]
	if !ok {
		return errors.New("beacon handler not found")
	}
	return beaconHandler.BeaconCommand(client, cmdline, beaconData, args)
}

func (h *Handler) BeaconPackData(beaconData response.BeaconData, tasks []response.TaskData) ([]byte, error) {
	beaconHandler, ok := h.BeaconHandlers[beaconData.Name]
	if !ok {
		return nil, errors.New("beacon handler not found")
	}
	return beaconHandler.BeaconPackData(beaconData, tasks)
}

func (h *Handler) BeaconProcessData(beaconData response.BeaconData, packedData []byte) ([]byte, error) {
	beaconHandler, ok := h.BeaconHandlers[beaconData.Name]
	if !ok {
		return nil, errors.New("module not found")
	}
	return beaconHandler.BeaconProcessData(beaconData, packedData)
}
