package model

type NetworkConfig struct {
	Network struct {
		Version   int                 `yaml:"version"`
		Ethernets map[string]Ethernet `yaml:"ethernets"`
	} `yaml:"network"`
}

type Ethernet struct {
	DHCP4       bool        `yaml:"dhcp4"`
	Nameservers Nameservers `yaml:"nameservers"`
}

type Nameservers struct {
	Addresses []string `yaml:"addresses"`
}

type NetworkConfigOption func(*NetworkConfig)

func WithEthernet(name string, ethernet Ethernet) NetworkConfigOption {
	return func(networkConfig *NetworkConfig) {
		networkConfig.Network.Ethernets[name] = ethernet
	}
}

func NewNetworkConfig(options ...NetworkConfigOption) *NetworkConfig {
	// Initialize the networkConfig struct
	networkConfig := &NetworkConfig{
		// Network: Network{
		// 	Version:   2,
		// 	Ethernets: make(map[string]Ethernet),
		// },
	}

	for _, option := range options {
		option(networkConfig)
	}

	return networkConfig
}
