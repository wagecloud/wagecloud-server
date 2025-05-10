package vm

import "github.com/wagecloud/wagecloud-server/internal/model"

type VMStatus string

const (
	VMStatusRunning VMStatus = "running"
	VMStatusStopped VMStatus = "stopped"
)

type VM struct {
	model.VM
	Status VMStatus `json:"status"`
}
