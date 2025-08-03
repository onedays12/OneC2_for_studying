package server

import (
	"OneServer/logs"
	"OneServer/utils/crypt"
	"OneServer/utils/response"
	"fmt"
	"go.uber.org/zap"
	"time"
)

func (ts *TeamServer) TaskCreate(beaconId string, cmdline string, client string, taskData response.TaskData) {
	value, ok := ts.Beacons.Get(beaconId)
	if !ok {
		logs.Logger.Error("TsTaskCreate: beacon not found", zap.String("beaconId", beaconId))
		return
	}

	beacon, _ := value.(*Beacon)
	if beacon.Active == false {
		return
	}

	if taskData.TaskId == "" {
		taskData.TaskId, _ = crypt.GenerateUID(8)
	}
	taskData.BeaconId = beaconId
	taskData.CommandLine = cmdline
	taskData.Client = client
	taskData.Computer = beacon.Data.Computer
	taskData.StartDate = time.Now().Unix()
	if taskData.Completed {
		taskData.FinishDate = taskData.StartDate
	}

	taskData.User = beacon.Data.Username
	if beacon.Data.Impersonated != "" {
		taskData.User += fmt.Sprintf(" [%s]", beacon.Data.Impersonated)
	}

	switch taskData.Type {

	case TYPE_TASK:
		beacon.TasksQueue.Put(taskData)
	default:
		break
	}
}

func (ts *TeamServer) TaskGetAll(beaconId string, availableSize int) ([]response.TaskData, error) {
	value, ok := ts.Beacons.Get(beaconId)
	if !ok {
		return nil, fmt.Errorf("TaskGetAll: beacon %v not found", beaconId)
	}
	beacon, _ := value.(*Beacon)

	var tasks []response.TaskData
	tasksSize := 0

	for i := uint(0); i < beacon.TasksQueue.Len(); i++ {
		value, ok = beacon.TasksQueue.Get(i)
		if ok {
			taskData := value.(response.TaskData)
			if tasksSize+len(taskData.Data) < availableSize {
				tasks = append(tasks, taskData)
				beacon.TasksQueue.Delete(i)
				i--
				tasksSize += len(taskData.Data)
			} else {
				break
			}
		} else {
			break
		}
	}

	return tasks, nil
}
