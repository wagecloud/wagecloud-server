package cloudinit

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
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

type CloudConfig struct {
	Users []UserDataParams `yaml:"users"` // represents list of users in user-data file
}

type UserDataParams struct {
	Name       string   `yaml:"name"`
	SSHKeys    []string `yaml:"ssh-authorized-keys"`
	Passwd     string   `yaml:"passwd,omitempty"`
	LockPasswd bool     `yaml:"lock_passwd"`
	Groups     string   `yaml:"groups,omitempty"`
	Shell      string   `yaml:"shell,omitempty"`
}

func WriteUserData(customer *Customer) error {
	fmt.Println("Name: ", customer.Name)
	fmt.Println("SSHKeys: ", customer.SSHKeys)
	fmt.Println("Passwd: ", customer.Passwd)
	fmt.Println("------")

	userData1 := &UserDataParams{
		Name:       customer.Name,
		SSHKeys:    customer.SSHKeys,
		Passwd:     customer.Passwd,
		LockPasswd: customer.LockPasswd,
		Groups:     "sudo",
		Shell:      "/bin/bash",
	}

	folderPath := fmt.Sprintf("./cloud-init-files/%s/%s", customer.Name, "vm1")
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

func GenerateUserData(params []UserDataParams) string {
	cloudConfig := &CloudConfig{
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
