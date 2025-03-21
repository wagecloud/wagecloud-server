package service

import (
	"github.com/wagecloud/wagecloud-server/internal/repository"
	"github.com/wagecloud/wagecloud-server/internal/service/cloudinit"
	"github.com/wagecloud/wagecloud-server/internal/service/libvirt"
	"github.com/wagecloud/wagecloud-server/internal/service/qemu"
)

type Service struct {
	Cloudinit *cloudinit.Service
	Libvirt   *libvirt.Service
	Qemu      *qemu.Service
}

func New(repo *repository.Repository) *Service {
	cloudinitService := cloudinit.NewService(repo)
	libvirtService := libvirt.NewService(repo)
	qemuService := qemu.NewService(repo)

	return &Service{
		Cloudinit: cloudinitService,
		Libvirt:   libvirtService,
		Qemu:      qemuService,
	}
}
