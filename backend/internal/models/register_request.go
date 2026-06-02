package models

type RegisterRequest struct {
	Username    string `json:"username"`
	MachineName string `json:"machine_name"`
	OSType      string `json:"os_type"`
}
