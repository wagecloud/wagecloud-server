package os

import (
	"context"

	"github.com/wagecloud/wagecloud-server/internal/model"
	"github.com/wagecloud/wagecloud-server/internal/repository"
	"github.com/wagecloud/wagecloud-server/internal/service/libvirt"
)


var _ ServiceInterface = (*Service)(nil)


type Service struct {
	repo repository.Repository
	libvirt libvirt.ServiceInterface
}

type ServiceInterface interface {

}
