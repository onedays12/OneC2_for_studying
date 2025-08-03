package response

type BeaconFile struct {
	FileName    string `json:"name"`
	FileContent string `json:"file_content"`
}

type BeaconData struct {
	Id           string `json:"b_id"`
	Name         string `json:"b_name"`
	SessionKey   []byte `json:"b_session_key"`
	Listener     string `json:"b_listener"`
	Async        bool   `json:"b_async"`
	ExternalIP   string `json:"b_external_ip"`
	InternalIP   string `json:"b_internal_ip"`
	GmtOffset    int    `json:"b_gmt_offset"`
	Sleep        uint   `json:"b_sleep"`
	Jitter       uint   `json:"b_jitter"`
	Pid          string `json:"b_pid"`
	Tid          string `json:"b_tid"`
	Arch         string `json:"b_arch"`
	Elevated     bool   `json:"b_elevated"`
	Process      string `json:"b_process"`
	Os           int    `json:"b_os"`
	OsDesc       string `json:"b_os_desc"`
	Domain       string `json:"b_domain"`
	Computer     string `json:"b_computer"`
	Username     string `json:"b_username"`
	Impersonated string `json:"b_impersonated"`
	OemCP        int    `json:"b_oemcp"`
	ACP          int    `json:"b_acp"`
	CreateTime   int64  `json:"b_create_time"`
	LastTick     int    `json:"b_last_tick"`
	KillDate     int    `json:"b_killdate"`
	WorkingTime  int    `json:"b_workingtime"`
	Tags         string `json:"b_tags"`
	Mark         string `json:"b_mark"`
	Color        string `json:"b_color"`
}
