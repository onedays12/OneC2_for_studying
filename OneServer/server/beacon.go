package server

import (
	"OneServer/utils/request"
	"OneServer/utils/response"
	"OneServer/utils/safeType"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"time"
)

func (ts *TeamServer) BeaconGenerate(beaconConfig request.BeaconConfig, listenerConfig request.ConfigDetail) ([]byte, string, error) {
	return ts.Handler.BeaconGenerate(beaconConfig, listenerConfig)
}

func (ts *TeamServer) BeaconIsExists(beaconId string) bool {
	return ts.Beacons.Contains(beaconId)
}

func (ts *TeamServer) BeaconCreate(beaconId string, beat []byte, listenerName string, ExternalIP string, Async bool) error {

	if len(beat) < 4 {
		return errors.New("data too short")
	}

	// 读取前4字节大端长度
	length := binary.BigEndian.Uint32(beat[:4])
	if int(length) > len(beat)-4 {
		return errors.New("invalid string length")
	}

	// 读取字符串内容（去掉末尾的 \x00）
	strBytes := beat[4 : 4+int(length)]
	beaconName := string(strBytes)
	if len(beaconName) > 0 && beaconName[len(beaconName)-1] == '\x00' {
		beaconName = beaconName[:len(beaconName)-1]
	}

	// 剩余数据
	restbeat := beat[4+int(length):]

	if restbeat == nil {
		return fmt.Errorf("beacon %v does not register", beaconId)
	}

	ok := ts.Beacons.Contains(beaconId)
	if ok {
		return fmt.Errorf("beacon %v already exists", beaconId)
	}

	beaconData, err := ts.Handler.BeaconCreate(beaconName, restbeat)
	if err != nil {
		return err
	}

	beaconData.Name = beaconName
	beaconData.Id = beaconId
	beaconData.Listener = listenerName
	beaconData.ExternalIP = ExternalIP
	beaconData.CreateTime = time.Now().Unix()
	beaconData.LastTick = int(time.Now().Unix())
	beaconData.Async = Async
	beaconData.Tags = ""
	beaconData.Mark = ""
	beaconData.Color = ""

	value, ok := ts.Listeners.Get(listenerName)
	if !ok {
		return fmt.Errorf("listener %v does not exists", listenerName)
	}

	lType := strings.Split(value.(response.ListenerData).Type, "/")[0]
	if lType == "internal" {
		beaconData.Mark = "Unlink"
	}

	beacon := &Beacon{
		Data:       beaconData,
		TasksQueue: safeType.NewSlice(),
		Active:     true,
	}

	ts.Beacons.Put(beaconData.Id, beacon)

	return nil
}

func (ts *TeamServer) BeaconCommand(beaconName string, beaconId string, clientName string, cmdline string, args map[string]any) error {
	value, ok := ts.Beacons.Get(beaconId)
	if !ok {
		return fmt.Errorf("beacon '%v' does not exist", beaconId)
	}
	beacon, _ := value.(*Beacon)

	if beacon.Active == false {
		return fmt.Errorf("beacon '%v' not active", beaconId)
	}

	return ts.Handler.BeaconCommand(clientName, cmdline, beaconName, beacon.Data, args)
}

func (ts *TeamServer) BeaconGetAllTasks(beaconId string, maxDataSize int) ([]byte, error) {

	value, ok := ts.Beacons.Get(beaconId)
	if !ok {
		return nil, fmt.Errorf("beacon type %v does not exists", beaconId)
	}

	beacon, _ := value.(*Beacon)

	tasksCount := beacon.TasksQueue.Len()
	if tasksCount > 0 {

		tasks, err := ts.TaskGetAll(beacon.Data.Id, maxDataSize)
		if err != nil {
			return nil, err
		}

		respData, err := ts.Handler.BeaconPackData(beacon.Data, tasks)
		if err != nil {
			return nil, err
		}

		return respData, nil
	}
	return nil, nil
}

func (ts *TeamServer) BeaconProcessData(beaconId string, bodyData []byte) error {
	value, ok := ts.Beacons.Get(beaconId)
	if !ok {
		return fmt.Errorf("beacon type %v does not exists", beaconId)
	}
	beacon, _ := value.(*Beacon)

	if beacon.Data.Async {
		beacon.Data.LastTick = int(time.Now().Unix())
		beacon.Tick = true
	}

	if beacon.Data.Mark == "Inactive" {
		beacon.Data.Mark = ""
	}

	if len(bodyData) > 4 {
		_, err := ts.Handler.BeaconProcessData(beacon.Data, bodyData)
		return err
	}

	return nil
}
