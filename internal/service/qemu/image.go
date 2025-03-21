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

func (s *Service) ImageCreate(baseImgPath string, cloneImgPath string) error {
	if !fileExists(baseImgPath) {
		return fmt.Errorf("base image not found")
	}

	cmd := exec.Command("qemu-img", // name
		"create",
		"-b",
		baseImgPath,
		"-f",
		"qcow2",
		"-F",
		"qcow2",
		cloneImgPath,
		"2G",
	)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create image: %s", err)
	}

	return nil
}

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
