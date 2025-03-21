package model

type Metadata struct {
	InstanceID    string `yaml:"instance-id"`
	LocalHostname string `yaml:"local-hostname"`
}

type MetadataOption func(*Metadata)

func WithMetadataInstanceID(instanceID string) MetadataOption {
	return func(metadata *Metadata) {
		metadata.InstanceID = instanceID
	}
}

func WithMetadataLocalHostname(localHostname string) MetadataOption {
	return func(metadata *Metadata) {
		metadata.LocalHostname = localHostname
	}
}

func NewMetadata(options ...MetadataOption) *Metadata {
	// Initialize the metadata struct
	metadata := &Metadata{}

	for _, option := range options {
		option(metadata)
	}

	return metadata
}
