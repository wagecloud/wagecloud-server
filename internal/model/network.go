package model

type Network struct {
	ID        string `json:"id"`
	PrivateIP string `json:"private_ip"`
	CreatedAt int64  `json:"created_at"`
}
