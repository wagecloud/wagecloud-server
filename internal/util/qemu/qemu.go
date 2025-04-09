package qemu

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/wagecloud/wagecloud-server/config"
)

func CreateImage(baseImgFile string, cloneImgFile string, size uint) error {
	if config.GetConfig().App.BaseImageDir == "" {
		return fmt.Errorf("base image dir not set")
	}

	baseImgPath := path.Join(
		config.GetConfig().App.BaseImageDir,
		baseImgFile,
	)

	if !exist(baseImgPath) {
		return fmt.Errorf("base image not found")
	}

	if config.GetConfig().App.ImageDir == "" {
		return fmt.Errorf("image dir not set")
	}

	if !exist(config.GetConfig().App.ImageDir) {
		os.MkdirAll(config.GetConfig().App.ImageDir, 0777)
	}

	cloneImgPath := path.Join(
		config.GetConfig().App.ImageDir,
		cloneImgFile,
	)

	sizeStr := fmt.Sprintf("%dG", size)

	// Eg: qemu-img create -b ubuntu_amd64.img -f qcow2 -F qcow2 ubuntu_amd64_mod.img 10G
	// set permissions to 777
	cmd := exec.Command("qemu-img",
		"create", "-b",
		baseImgPath,
		"-f",
		"qcow2",
		"-F",
		"qcow2",
		cloneImgPath,
		// "10G", // TODO: add volume params
		sizeStr, // G for GB
	)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create image: %s", err)
	}

	return nil
}

// neet account id, base path with clone path

func CreateImageWithPath(basePath string, clonePath string, size uint) error {
	if !exist(basePath) {
		return fmt.Errorf("base image not found")
	}

	if err := os.MkdirAll(path.Dir(clonePath), 0777); err != nil {
		return fmt.Errorf("failed to create directory: %s", err)
	}

	sizeStr := fmt.Sprintf("%dG", size)
	fmt.Println("sizeStr", sizeStr)
	//eg: sudo qemu-img create -b /var/lib/libvirt/images/alexng/base/focal-server-cloudimg-amd64.img
	// -f qcow2 -F qcow2
	// /var/lib/libvirt/images/alexng/7a4a5c55-000c-44d5-b41e-903b71bf32fe/focal-server-cloudimg-amd64.img
	// set permissions to 777
	cmd := exec.Command("qemu-img",
		"create", "-b",
		basePath,
		"-f",
		"qcow2",
		"-F",
		"qcow2",
		clonePath,
		// "10G", // TODO: add volume params
		sizeStr, // G for GB
	)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create image: %s", err)
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

func exist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
