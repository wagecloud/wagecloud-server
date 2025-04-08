package service

import (
	"github.com/wagecloud/wagecloud-server/internal/repository"
	"github.com/wagecloud/wagecloud-server/internal/service/account"
	"github.com/wagecloud/wagecloud-server/internal/service/cloudinit"
	"github.com/wagecloud/wagecloud-server/internal/service/libvirt"
	"github.com/wagecloud/wagecloud-server/internal/service/qemu"
)

type Service struct {
	Account   account.ServiceInterface
	Cloudinit cloudinit.ServiceInterface
	Libvirt   libvirt.ServiceInterface
	Qemu      qemu.ServiceInterface
}

func New(repo *repository.RepositoryImpl) *Service {
	accountSvc := account.NewService(repo)
	cloudinitSvc := cloudinit.NewService(repo)
	libvirtSvc := libvirt.NewService(repo)
	qemuSvc := qemu.NewService(repo)

	return &Service{
		Account:   accountSvc,
		Cloudinit: cloudinitSvc,
		Libvirt:   libvirtSvc,
		Qemu:      qemuSvc,
	}
}
