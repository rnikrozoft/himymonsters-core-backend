package model

type RegisterWithDeviceID struct {
	Username string `json:"username"`
	Created  bool   `json:"created"`
}

type Wallet struct {
	Coin   int64 `json:"coin"`
	Dimond int64 `json:"dimond"`
}
