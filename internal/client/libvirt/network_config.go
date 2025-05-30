package libvirt

type NetworkConfig struct {
	Network Network `json:"network" yaml:"network"`
}
type Match struct {
	Driver string `json:"driver,omitempty" yaml:"driver,omitempty"`
}
type Ethernet struct {
	Match   Match  `json:"match" yaml:"match"`
	Dhcp4   bool   `json:"dhcp4,omitempty" yaml:"dhcp4,omitempty"`
	SetName string `json:"set-name,omitempty" yaml:"set-name,omitempty"`
}
type Network struct {
	Version   int                 `json:"version,omitempty" yaml:"version,omitempty"`
	Ethernets map[string]Ethernet `json:"ethernets,omitempty" yaml:"ethernets,omitempty"`
}

func NewDefaultNetworkConfig() NetworkConfig {
	n := NetworkConfig{
		Network: Network{
			Version:   2,
			Ethernets: make(map[string]Ethernet),
		},
	}

	n.Network.Ethernets["eth0"] = Ethernet{
		Match: Match{
			Driver: "virtio_net",
		},
		Dhcp4:   true,
		SetName: "eth0",
	}

	return n
}
