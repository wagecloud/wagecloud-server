package libvirt

type Userdata struct {
	Users  []User   `json:"users,omitempty" yaml:"users,omitempty"`
	Runcmd []string `json:"runcmd,omitempty" yaml:"runcmd,omitempty"`
}

type User struct {
	Name              string   `json:"name,omitempty" yaml:"name,omitempty"`
	SSHAuthorizedKeys []string `json:"ssh-authorized-keys,omitempty" yaml:"ssh-authorized-keys,omitempty"`
	Passwd            string   `json:"passwd,omitempty" yaml:"passwd,omitempty"`
	LockPasswd        bool     `json:"lock_passwd" yaml:"lock_passwd"`
	Groups            string   `json:"groups,omitempty" yaml:"groups,omitempty"`
	Sudo              string   `json:"sudo,omitempty" yaml:"sudo,omitempty"`
	Shell             string   `json:"shell,omitempty" yaml:"shell,omitempty"`
}

func NewDefaultUserdata() Userdata {
	u := Userdata{
		Users: []User{
			NewDefaultUser(),
		},
		Runcmd: NewDefaultRuncmd(),
	}
	return u
}

func NewDefaultUser() User {
	u := User{
		Name:              "root",
		SSHAuthorizedKeys: []string{},
		Passwd:            "",
		LockPasswd:        false,
		Groups:            "sudo",
		Sudo:              "ALL=(ALL) NOPASSWD:ALL",
		Shell:             "/bin/bash",
	}
	return u
}

func NewDefaultRuncmd() []string {
	return []string{
		// "rm -f /etc/machine-id",
		// "systemd-machine-id-setup",
		// "ln -s /etc/machine-id /var/lib/dbus/machine-id",
	}
}
