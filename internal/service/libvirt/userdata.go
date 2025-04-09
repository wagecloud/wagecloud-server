package libvirt

type Userdata struct {
	Users []User `json:"users,omitempty" yaml:"users,omitempty"`
}

type User struct {
	Name              string   `json:"name,omitempty" yaml:"name,omitempty"`
	SSHAuthorizedKeys []string `json:"ssh-authorized-keys,omitempty" yaml:"ssh-authorized-keys,omitempty"`
	Passwd            string   `json:"passwd,omitempty" yaml:"passwd,omitempty"`
	LockPasswd        bool     `json:"lock_passwd,omitempty" yaml:"lock_passwd,omitempty"`
	Groups            string   `json:"groups,omitempty" yaml:"groups,omitempty"`
	Sudo              string   `json:"sudo,omitempty" yaml:"sudo,omitempty"`
	Shell             string   `json:"shell,omitempty" yaml:"shell,omitempty"`
}

func NewDefaultUserdata() Userdata {
	u := Userdata{
		Users: []User{
			NewDefaultUser(),
		},
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
