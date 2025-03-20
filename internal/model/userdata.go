package model

type Userdata struct {
	Name       string   `yaml:"name"`
	SSHKeys    []string `yaml:"ssh-authorized-keys"`
	Passwd     string   `yaml:"passwd,omitempty"`
	LockPasswd bool     `yaml:"lock_passwd"`
	Groups     string   `yaml:"groups,omitempty"`
	Sudo       string   `yaml:"sudo,omitempty"`
	Shell      string   `yaml:"shell,omitempty"`
}
