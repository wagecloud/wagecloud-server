package service

import (
	"github.com/wagecloud/wagecloud-server/internal/repository"
	"github.com/wagecloud/wagecloud-server/internal/service/account"
	"github.com/wagecloud/wagecloud-server/internal/service/arch"
	"github.com/wagecloud/wagecloud-server/internal/service/libvirt"
	"github.com/wagecloud/wagecloud-server/internal/service/network"
	"github.com/wagecloud/wagecloud-server/internal/service/os"
	"github.com/wagecloud/wagecloud-server/internal/service/vm"
)

type Service struct {
	Account account.ServiceInterface
	Arch    arch.ServiceInterface
	Libvirt libvirt.ServiceInterface
	VM      vm.ServiceInterface
	OS      os.ServiceInterface
	Network network.ServiceInterface
}

func New(repo *repository.RepositoryImpl) *Service {
	accountSvc := account.NewService(repo)
	archSvc := arch.NewService(repo)
	libvirtSvc := libvirt.NewService(repo)
	vmSvc := vm.NewService(repo, libvirtSvc)
	osSvc := os.NewService(repo, libvirtSvc)
	networkSvc := network.NewService(repo)
	return &Service{
		Account: accountSvc,
		Arch:    archSvc,
		Libvirt: libvirtSvc,
		VM:      vmSvc,
		OS:      osSvc,
		Network: networkSvc,
	}
}
