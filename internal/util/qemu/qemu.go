package qemu

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/wagecloud/wagecloud-server/config"
	"github.com/wagecloud/wagecloud-server/internal/util/file"
	"github.com/wagecloud/wagecloud-server/internal/util/transaction"
)

func CreateImage(baseImagePath string, cloneImagePath string, size uint) error {
	if config.GetConfig().App.BaseImageDir == "" {
		return fmt.Errorf("base image dir not set")
	}

	if !file.Exists(baseImagePath) {
		return fmt.Errorf("base image not found")
	}

	cloneImgPath := path.Join(
		config.GetConfig().App.VMImageDir,
		cloneImagePath,
	)

	sizeStr := fmt.Sprintf("%dG", size)

	tx := transaction.NewTransaction()

	tx.Add(func() error {
		// Eg: qemu-img create -b ubuntu_amd64.img -f qcow2 -F qcow2 ubuntu_amd64_mod.img 10G
		// set permissions to 777
		cmd := exec.Command("qemu-img",
			"create", "-b",
			baseImagePath,
			"-f",
			"qcow2",
			"-F",
			"qcow2",
			cloneImagePath,
			sizeStr, // G for GB
		)
		return cmd.Run()
	}, func() error {
		return os.Remove(cloneImgPath)
	})

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %s", err)
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
