package vm

import "github.com/wagecloud/wagecloud-server/internal/repository"

var _ ServiceInterface = (*Service)(nil)

type Service struct {
	repo repository.Repository
}

type ServiceInterface interface {
	CreateVM() error
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateVM() error {
	// Implementation for creating a VM goes here.
	return nil
}
