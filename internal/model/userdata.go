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

type UserdataOption func(*Userdata)

func WithUserdataName(name string) UserdataOption {
	return func(userdata *Userdata) {
		userdata.Name = name
	}
}

func WithUserdataSSHKeys(sshKeys []string) UserdataOption {
	return func(userdata *Userdata) {
		userdata.SSHKeys = sshKeys
	}
}

func WithUserdataPasswd(passwd string) UserdataOption {
	return func(userdata *Userdata) {
		userdata.Passwd = passwd
	}
}

func WithUserdataLockPasswd(lockPasswd bool) UserdataOption {
	return func(userdata *Userdata) {
		userdata.LockPasswd = lockPasswd
	}
}

func WithUserdataGroups(groups string) UserdataOption {
	return func(userdata *Userdata) {
		userdata.Groups = groups
	}
}

func WithUserdataSudo(sudo string) UserdataOption {
	return func(userdata *Userdata) {
		userdata.Sudo = sudo
	}
}

func WithUserdataShell(shell string) UserdataOption {
	return func(userdata *Userdata) {
		userdata.Shell = shell
	}
}

func NewUserdata(options ...UserdataOption) *Userdata {
	// Initialize the userdata struct
	userdata := &Userdata{
		Groups: "sudo",
		Sudo:   "ALL=(ALL) NOPASSWD:ALL",
		Shell:  "/bin/bash",
	}

	for _, option := range options {
		option(userdata)
	}

	return userdata
}
