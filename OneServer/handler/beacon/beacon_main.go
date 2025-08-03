package beacon

import (
	"OneServer/utils/request"
	"OneServer/utils/response"
	"errors"
	"time"
)

var (
	beaconHandler *BeaconHandler
)

func NewBeaconHandler(ts any) any {
	beaconHandler = &BeaconHandler{
		ts: ts.(TeamServer),
	}
	return beaconHandler
}

func (b *BeaconHandler) BeaconGenerate(beaconConfig request.GenerateConfig, listenerConfig request.ConfigDetail) ([]byte, string, error) {

	beaconProfile, err := b.BeaconGenerateProfile(beaconConfig, listenerConfig)
	if err != nil {
		return nil, "", err
	}
	return b.BeaconBuild(beaconProfile, beaconConfig, listenerConfig)
}

func (b *BeaconHandler) BeaconCreate(beat []byte) (response.BeaconData, error) {
	return b.CreateBeacon(beat)
}

func (b *BeaconHandler) BeaconCommand(client string, cmdline string, beaconData response.BeaconData, args map[string]any) error {
	command, ok := args["command"].(string)
	if !ok {
		return errors.New("'command' must be set")
	}

	taskData, err := b.CreateTask(beaconData, command, args)
	if err != nil {
		return err
	}

	b.ts.TaskCreate(beaconData.Id, cmdline, client, taskData)

	return nil
}

func (b *BeaconHandler) BeaconPackData(beaconData response.BeaconData, tasks []response.TaskData) ([]byte, error) {
	packedData, err := b.PackTasks(tasks)
	if err != nil {
		return nil, err
	}

	return b.EncryptData(packedData, beaconData.SessionKey)
}

func (b *BeaconHandler) BeaconProcessData(beaconData response.BeaconData, packedData []byte) ([]byte, error) {
	decryptData, err := b.DecryptData(packedData, beaconData.SessionKey)
	if err != nil {
		return nil, err
	}

	taskData := response.TaskData{
		Type:        TYPE_TASK,
		BeaconId:    beaconData.Id,
		FinishDate:  time.Now().Unix(),
		MessageType: MESSAGE_SUCCESS,
		Completed:   true,
		Sync:        true,
	}

	err = b.ProcessTasksResult(b.ts, beaconData, taskData, decryptData)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
