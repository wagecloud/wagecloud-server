package qemu

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/wagecloud/wagecloud-server/internal/repository"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateImage(baseImgPath string, cloneImgPath string) error {
	if !fileExists(baseImgPath) {
		return fmt.Errorf("base image not found")
	}

	// Eg: qemu-img create -b ubuntu_amd64.img -f qcow2 -F qcow2 ubuntu_amd64_mod.img 10G
	cmd := exec.Command("qemu-img",
		"create",
		"-b",
		baseImgPath,
		"-f",
		"qcow2",
		"-F",
		"qcow2",
		cloneImgPath,
		"10G", // TODO: add volumn params
	)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create image: %s", err)
	}

	return nil
}

// func (s *Service) Convert(imgPath string, format string, destPath string) error {
// 	if !fileExists(imgPath) {
// 		return fmt.Errorf("image not found")
// 	}

// 	// Eg. qemu-img convert -f qcow2 -O raw focal-server-cloudimg-amd64.img focal-server-cloudimg-amd64.raw
// 	cmd := exec.Command("qemu-img", // name
// 		"convert",
// 		"-f",
// 		"qcow2",
// 		"-O",
// 		format,
// 		imgPath,
// 		destPath,
// 	)

// 	if err := cmd.Run(); err != nil {
// 		return fmt.Errorf("failed to convert image: %s", err)
// 	}

// 	return nil
// }

// func ImageResize(imgPath string, vol *Volumn) error {
// 	if !fileExists(imgPath) {
// 		return fmt.Errorf("image not found")
// 	}

// 	cmd := exec.Command("qemu-img", // name
// 		"resize",
// 		imgPath,
// 		fmt.Sprintf("%d%s", vol.size, vol.unit),
// 	)

// 	err := cmd.Run()
// 	if err != nil {
// 		return fmt.Errorf("failed to resize image: %s", err)
// 	}

// 	return nil
// }

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
