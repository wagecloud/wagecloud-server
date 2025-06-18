package instancemodel

import "time"

type Status string
type LogType string

const (
	StatusUnknown Status = "STATUS_UNKNOWN"
	StatusPending Status = "STATUS_PENDING"
	StatusRunning Status = "STATUS_RUNNING"
	StatusStopped Status = "STATUS_STOPPED"
	StatusError   Status = "STATUS_ERROR"

	LogUnknown LogType = "LOG_TYPE_UNKNOWN"
	LogInfo    LogType = "LOG_TYPE_INFO"
	LogWarning LogType = "LOG_TYPE_WARNING"
	LogError   LogType = "LOG_TYPE_ERROR"
)

type Instance struct {
	ID        string    `json:"id"`
	AccountID int64     `json:"account_id"`
	OSID      string    `json:"os_id"`
	ArchID    string    `json:"arch_id"`
	RegionID  string    `json:"region_id"`
	Name      string    `json:"name"`
	CPU       int32     `json:"cpu"`
	RAM       int32     `json:"ram"`     // in MB
	Storage   int32     `json:"storage"` // in GB
	CreatedAt time.Time `json:"created_at"`
}

type InstanceMonitor struct {
	ID           string  `json:"id"`
	Status       Status  `json:"status"`
	CPUUsage     float64 `json:"cpu_usage"`     // in percentage
	RAMUsage     float64 `json:"ram_usage"`     // in MB
	StorageUsage float64 `json:"storage_usage"` // in MB
	NetworkIn    float64 `json:"network_in"`    // in MB
	NetworkOut   float64 `json:"network_out"`   // in MB
}

type Network struct {
	ID         int64   `json:"id"`
	InstanceID string  `json:"instance_id"`
	PrivateIP  string  `json:"private_ip"`
	MacAddress string  `json:"mac_address,omitempty"`
	PublicIP   *string `json:"public_ip"`
}

type Domain struct {
	ID        int64  `json:"id"`
	NetworkID int64  `json:"network_id"`
	Name      string `json:"name"`
}

type InstanceLog struct {
	ID          int64   `json:"id"`
	InstanceID  string  `json:"instance_id"`
	Type        LogType `json:"type"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	CreatedAt   time.Time
}

type Region struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
