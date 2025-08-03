package request

type BeaconConfig struct {
	ListenerName string         `json:"listener_name"`
	ListenerType string         `json:"listener_type"`
	BeaconType   string         `json:"beacon_type"`
	Config       GenerateConfig `json:"config"`
}

type GenerateConfig struct {
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
}

type CommandData struct {
	BeaconType string         `json:"beacon_type"`
	BeaconId   string         `json:"beacon_id"`
	CmdLine    string         `json:"cmdline"`
	Data       map[string]any `json:"data"`
}
