package cloudinit

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"

	"gopkg.in/yaml.v3"
)

type HashType string

const (
	SHA256 HashType = "sha256"
	SHA512 HashType = "sha512"
)

func NewHashFunc(hashType HashType) hash.Hash {
	switch hashType {
	case SHA256:
		return sha256.New()
	default:
		return sha512.New()
	}
}

type CloudConfig struct {
	Users []UserDataParams `yaml:"users"`
}

type UserDataParams struct {
	Name       string   `yaml:"name"`
	SSHKeys    []string `yaml:"ssh-authorized-keys"`
	Passwd     string   `yaml:"passwd,omitempty"`
	LockPasswd bool     `yaml:"lock_passwd,omitempty"`
	Groups     string   `yaml:"groups,omitempty"`
	Shell      string   `yaml:"shell,omitempty"`
}

func GenerateUserData(params *UserDataParams) string {
	cloudConfig := &CloudConfig{
		Users: []UserDataParams{*params, *params},
	}

	yamlData, err := yaml.Marshal(cloudConfig)

	if err != nil {
		panic(err)
	}

	return string(yamlData)
}

func HashPassword(plainPassword string, hashType HashType) (string, error) {
	if plainPassword == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	hashFunc := NewHashFunc(hashType)
	hashFunc.Write([]byte(plainPassword))
	return hex.EncodeToString(hashFunc.Sum(nil)), nil
}
