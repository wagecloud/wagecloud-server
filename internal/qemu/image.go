package qemu

import (
	"fmt"
	"os"
	"os/exec"
)

type OSBaseImage string
type VolumnUnit string

const (
	UBUNTU OSBaseImage = "focal-server-cloudimg-amd64.img"
	DEBIAN             = "ahihi"
)

const (
	G VolumnUnit = "G"
)

type Volumn struct {
	size uint
	unit VolumnUnit
}

func ImageCreate(baseImgPath string, cloneImgPath string, vol *Volumn) error {
	if !fileExists(string(UBUNTU)) {
		return fmt.Errorf("base image not found")
	}

	cmd := exec.Command("qemu-img", // name
		"create",
		"-b",
		string(UBUNTU),
		"-f",
		"qcow2",
		"-F",
		"qcow2",
		cloneImgPath,
		fmt.Sprintf("%d%s", vol.size, vol.unit),
	)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to create image: %s", err)
	}

	return nil
}

func ImageResize(imgPath string, vol *Volumn) error {
	if !fileExists(imgPath) {
		return fmt.Errorf("image not found")
	}

	cmd := exec.Command("qemu-img", // name
		"resize",
		imgPath,
		fmt.Sprintf("%d%s", vol.size, vol.unit),
	)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to resize image: %s", err)
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
