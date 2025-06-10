package instancemodel

type Status string

const (
	StatusPending    Status = "PENDING"
	StatusRunning    Status = "RUNNING"
	StatusStopped    Status = "STOPPED"
	StatusFailed     Status = "FAILED"
	StatusTerminated Status = "TERMINATED"
)

type Instance struct {
	ID        string `json:"id"`
	AccountID int64  `json:"account_id"`
	NetworkID string `json:"network_id"`
	OSID      string `json:"os_id"`
	ArchID    string `json:"arch_id"`
	Name      string `json:"name"`
	CPU       int32  `json:"cpu"`
	RAM       int32  `json:"ram"`     // in MB
	Storage   int32  `json:"storage"` // in GB
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type Network struct {
	ID        string `json:"id"`
	PrivateIP string `json:"private_ip"`
	CreatedAt int64  `json:"created_at"`
}

type OS struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
}

type Arch struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
}
