package libvirt

import "github.com/google/uuid"

type Metadata struct {
	InstanceID    string `json:"instance-id" yaml:"instance-id"`
	LocalHostname string `json:"local-hostname" yaml:"local-hostname"`
}

func NewDefaultMetadata() Metadata {
	id := uuid.New().String()

	return Metadata{
		InstanceID:    id,
		LocalHostname: id,
	}
}
