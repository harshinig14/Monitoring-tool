package dto

type DeviceResponse struct {
	UserID      int    `json:"user_id"`
	MachineName string `json:"machine_name"`
	OSType      string `json:"os_type"`
	Status      string `json:"status"`
	LastSeen    string `json:"last_seen"`
}
