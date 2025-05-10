package model

type VM struct {
	ID        string `json:"id"`
	AccountID int64  `json:"account_id"`
	NetworkID string `json:"network_id"`
	OsID      string `json:"os_id"`
	ArchID    string `json:"arch_id"`
	Name      string `json:"name"`
	Cpu       int32  `json:"cpu"`
	Ram       int32  `json:"ram"`     // In MB
	Storage   int32  `json:"storage"` // In GB
	CreatedAt int64  `json:"created_at"`
}
