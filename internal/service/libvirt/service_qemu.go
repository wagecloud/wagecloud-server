package libvirt

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/wagecloud/wagecloud-server/config"
	"github.com/wagecloud/wagecloud-server/internal/util/file"
)

type CreateImageParams struct {
	BaseImagePath  string
	CloneImagePath string
	Size           uint
}

func (s *Service) CreateImage(params CreateImageParams) error {
	if config.GetConfig().App.BaseImageDir == "" {
		return fmt.Errorf("base image dir not set")
	}

	if !file.Exists(params.BaseImagePath) {
		return fmt.Errorf("base image not found")
	}

	cloneImgPath := path.Join(
		config.GetConfig().App.VMImageDir,
		params.CloneImagePath,
	)

	sizeStr := fmt.Sprintf("%dG", params.Size)

	if err := s.tx.Do(func() error {
		// Eg: qemu-img create -b ubuntu_amd64.img -f qcow2 -F qcow2 ubuntu_amd64_mod.img 10G
		// set permissions to 777
		cmd := exec.Command("qemu-img",
			"create", "-b",
			params.BaseImagePath,
			"-f",
			"qcow2",
			"-F",
			"qcow2",
			params.CloneImagePath,
			sizeStr, // G for GB
		)
		return cmd.Run()
	}, func() error {
		return os.Remove(cloneImgPath)
	}); err != nil {
		return err
	}

	return nil
}

// func (s *Service) Convert(imgPath string, format string, destPath string) error {
// 	if !exist(imgPath) {
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
// 	if !exist(imgPath) {
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
