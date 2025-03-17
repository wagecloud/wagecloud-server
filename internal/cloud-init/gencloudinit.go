package cloudinit

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v3"
)

type HashType string

const (
	BCrypt HashType = "Bcrypt"
)

type Customer struct {
	Name       string
	SSHKeys    []string
	Passwd     string
	LockPasswd bool
}

type UserCloudConfig struct {
	Users []UserDataParams `yaml:"users"` // represents list of users in user-data file
}

type MetadataCloudConfig struct {
	InstanceID    string `yaml:"instance-id"`
	LocalHostname string `yaml:"local-hostname"`
}

type UserDataParams struct {
	Name       string   `yaml:"name"`
	SSHKeys    []string `yaml:"ssh-authorized-keys"`
	Passwd     string   `yaml:"passwd,omitempty"`
	LockPasswd bool     `yaml:"lock_passwd"`
	Groups     string   `yaml:"groups,omitempty"`
	Sudo       string   `yaml:"sudo,omitempty"`
	Shell      string   `yaml:"shell,omitempty"`
}

func WriteCloudInitFiles(customer *Customer, path string) error {
	err := WriteUserData(customer, path)
	if err != nil {
		return fmt.Errorf("failed to write user data: %s", err)
	}

	err = WriteMetaData(customer, path)
	if err != nil {
		return fmt.Errorf("failed to write metadata: %s", err)
	}

	return nil
}

func WriteUserData(customer *Customer, folderPath string) error {
	userData1 := &UserDataParams{
		Name:       customer.Name,
		SSHKeys:    customer.SSHKeys,
		Passwd:     customer.Passwd,
		LockPasswd: customer.LockPasswd,
		Sudo:       "ALL=(ALL) NOPASSWD:ALL",
		Groups:     "sudo",
		Shell:      "/bin/bash",
	}

	filepath := filepath.Join(folderPath, "user-data")

	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create folder: %s", err)
	}

	content := GenerateUserData([]UserDataParams{*userData1})
	os.WriteFile(filepath, []byte(content), 0644)
	fmt.Println(content)

	return nil
}

func WriteMetaData(customer *Customer, folderPath string) error {
	metaData := &MetadataCloudConfig{
		InstanceID:    customer.Name + "-" + genUUID(),
		LocalHostname: customer.Name,
	}

	filepath := filepath.Join(folderPath, "meta-data")

	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create folder: %s", err)
	}

	content, err := yaml.Marshal(metaData)

	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %s", err)
	}

	os.WriteFile(filepath, content, 0644)
	return nil
}

func GenerateUserData(params []UserDataParams) string {
	cloudConfig := &UserCloudConfig{
		Users: params,
	}

	yamlData, err := yaml.Marshal(cloudConfig)
	if err != nil {
		panic(err)
	}
	dataStr := string(yamlData)
	return fmt.Sprintf("#cloud-config\n%s", dataStr)
}

func HashPassword(plainPassword string, hashType HashType) (string, error) {
	if plainPassword == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	if hashType != BCrypt { // temp condition
		return "", fmt.Errorf("unsupported hash type")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %s", err)
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(plainPassword))

	if err != nil {
		return "", fmt.Errorf("failed to compare password: %s", err)
	}

	return string(hashedPassword), nil
}

func genUUID() string {
	return uuid.New().String()
}
