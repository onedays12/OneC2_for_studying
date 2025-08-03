package response

type TaskData struct {
	Type        int    `json:"t_type"`
	TaskId      string `json:"t_task_id"`
	BeaconId    string `json:"t_beacon_id"`
	Client      string `json:"t_client"`
	User        string `json:"t_user"`
	Computer    string `json:"t_computer"`
	StartDate   int64  `json:"t_start_date"`
	FinishDate  int64  `json:"t_finish_date"`
	Data        []byte `json:"t_data"`
	CommandLine string `json:"t_command_line"`
	MessageType int    `json:"t_message_type"`
	Message     string `json:"t_message"`
	ClearText   string `json:"t_clear_text"`
	Completed   bool   `json:"t_completed"`
	Sync        bool   `json:"t_sync"`
}
