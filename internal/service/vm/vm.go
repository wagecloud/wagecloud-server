package vm

import "github.com/wagecloud/wagecloud-server/internal/repository"

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateVM() error {
	// Implementation for creating a VM goes here.
	return nil
}
