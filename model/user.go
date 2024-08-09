package model

type RegisterWithDeviceID struct {
	DeviceID string `json:"device_id"`
	Username string `json:"username"`
}
