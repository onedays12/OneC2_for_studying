package request

type ListenerConfig struct {
	ListenerName string       `json:"name"`
	ConfigType   string       `json:"type"`
	Config       ConfigDetail `json:"config"`
}

type ConfigDetail struct {
	HostBind          string            `json:"host_bind"`
	PortBind          int               `json:"port_bind"`
	CallbackAddresses []string          `json:"callback_addresses"`
	SSL               bool              `json:"ssl"`
	SSLCert           []byte            `json:"ssl_cert"`
	SSLKey            []byte            `json:"ssl_key"`
	SSLCertPath       string            `json:"ssl_cert_path"`
	SSLKeyPath        string            `json:"ssl_key_path"`
	URI               string            `json:"uri"`
	HBHeader          string            `json:"hb_header"`
	HBPrefix          string            `json:"hb_prefix"`
	UserAgent         string            `json:"user_agent"`
	HostHeader        string            `json:"host_header"`
	RequestHeaders    map[string]string `json:"request_headers"`
	ResponseHeaders   map[string]string `json:"response_headers"`
	XForwardedFor     bool              `json:"x_forwarded_for"`
	PageError         string            `json:"page_error"`
	PagePayload       string            `json:"page_payload"`
	ServerHeaders     map[string]string `json:"server_headers"`
	Protocol          string            `json:"protocol"`
	EncryptKey        string            `json:"encrypt_key"`
}
