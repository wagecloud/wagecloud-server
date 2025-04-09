package service

import (
	"github.com/wagecloud/wagecloud-server/internal/repository"
	"github.com/wagecloud/wagecloud-server/internal/service/account"
	"github.com/wagecloud/wagecloud-server/internal/service/libvirt"
	"github.com/wagecloud/wagecloud-server/internal/service/qemu"
	"github.com/wagecloud/wagecloud-server/internal/service/vm"
)

type Service struct {
	Account account.ServiceInterface
	Libvirt libvirt.ServiceInterface
	Qemu    qemu.ServiceInterface
	VM      vm.ServiceInterface
}

func New(repo *repository.RepositoryImpl) *Service {
	accountSvc := account.NewService(repo)
	libvirtSvc := libvirt.NewService(repo)
	qemuSvc := qemu.NewService(repo)
	vmSvc := vm.NewService(repo, libvirtSvc)

	return &Service{
		Account: accountSvc,
		Libvirt: libvirtSvc,
		Qemu:    qemuSvc,
		VM:      vmSvc,
	}
}
