package profile

type Profile struct {
	TeamServerConfig     TeamServerConfig     `json:"TeamServer"`
	ServerResponseConfig ServerResponseConfig `json:"ServerResponse"`
	ZapConfig            ZapConfig            `json:"Zap"`
}

type TeamServerConfig struct {
	Host                  string   `json:"host"`
	Port                  int      `json:"port"`
	Endpoint              string   `json:"endpoint"`
	Password              string   `json:"password"`
	Cert                  string   `json:"cert"`
	Key                   string   `json:"key"`
	Extenders             []string `json:"extenders"`
	AccessTokenLiveHours  int      `json:"access_token_live_hours"`
	RefreshTokenLiveHours int      `json:"refresh_token_live_hours"`
	Env                   string   `json:"env"`
}

type ServerResponseConfig struct {
	Status   int               `json:"status"`
	Headers  map[string]string `json:"headers"`
	PagePath string            `json:"pagepath"`
	Page     string            `json:"-"`
}

type ZapConfig struct {
	Level          string `json:"level"`
	MaxSize        int    `json:"max_size"`
	MaxBackups     int    `json:"max_backups"`
	MaxAge         int    `json:"max_age"`
	IsConsolePrint bool   `json:"is_console_print"`
	Path           string `json:"path"`
}
