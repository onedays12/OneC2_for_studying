package response

import "OneServer/utils/request"

type ListenerData struct {
	Name       string               `json:"listener_name"`
	Type       string               `json:"listener_type"`
	BindHost   string               `json:"listener_bind_host"`
	BindPort   string               `json:"listener_bind_port"`
	BeaconAddr string               `json:"listener_beacon_addr"`
	Status     string               `json:"listener_status"`
	Data       request.ConfigDetail `json:"listener_data"`
}
