package libvirt

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/wagecloud/wagecloud-server/internal/logger"
	"github.com/wagecloud/wagecloud-server/internal/utils/saga"
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)

type ClientImpl struct {
	connect *libvirt.Connect
	saga    *saga.Saga
}

type Client interface {
	// CLOUDINIT
	CreateCloudinit(ctx context.Context, params CreateCloudinitParams) error
	CreateCloudinitByReader(ctx context.Context, params CreateCloudinitByReaderParams) error
	WriteCloudinit(ctx context.Context, userdata io.Reader, metadata io.Reader, networkConfig io.Reader, cloudinitFile io.Writer) error

	// DOMAIN
	GetDomain(ctx context.Context, domainID string) (Domain, error)
	ListDomains(ctx context.Context, params ListDomainsParams) ([]Domain, error)
	CreateDomain(ctx context.Context, domain Domain) error
	UpdateDomain(ctx context.Context, domainID string, params UpdateDomainParams) error
	DeleteDomain(ctx context.Context, domainID string) error
	StartDomain(ctx context.Context, domainID string) error
	StopDomain(ctx context.Context, domainID string) error

	// QEMU
	CreateImage(ctx context.Context, params CreateImageParams) error
}

const (
	QemuConnect = "qemu:///system"
)

var (
	ErrDomainNotFound = errors.New("domain not found")
)

func NewClient() Client {
	return &ClientImpl{
		connect: nil,
		saga:    saga.New(),
	}
}

func (s *ClientImpl) getConnect() (*libvirt.Connect, error) {
	if s.connect == nil {
		conn, err := libvirt.NewConnect(QemuConnect)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to libvirt: %v", err)
		}
		s.connect = conn
	}
	return s.connect, nil
}

func (s *ClientImpl) getDomain(domainID string) (*libvirt.Domain, error) {
	conn, err := s.getConnect()
	if err != nil {
		return nil, err
	}

	domain, err := conn.LookupDomainByUUIDString(domainID)
	if err != nil {
		return nil, ErrDomainNotFound
	}

	return domain, nil
}

func (s *ClientImpl) GetDomain(ctx context.Context, domainID string) (Domain, error) {
	conn, err := s.getConnect()
	if err != nil {
		return Domain{}, err
	}

	domain, err := conn.LookupDomainByUUIDString(domainID)
	if err != nil {
		return Domain{}, ErrDomainNotFound
	}

	return FromLibvirtToDomain(*domain)
}

type ListDomainsParams struct {
	Flags libvirt.ConnectListAllDomainsFlags
}

func (s *ClientImpl) ListDomains(ctx context.Context, params ListDomainsParams) ([]Domain, error) {
	conn, err := s.getConnect()
	if err != nil {
		return nil, err
	}

	domains, err := conn.ListAllDomains(params.Flags)
	if err != nil {
		return nil, fmt.Errorf("failed to list domains: %v", err)
	}

	domainsModel := make([]Domain, len(domains))

	for i, domain := range domains {
		model, err := FromLibvirtToDomain(domain)
		if err != nil {
			return nil, fmt.Errorf("failed to convert domain to model: %v", err)
		}

		domainsModel[i] = model
	}

	return domainsModel, nil
}

// CreateDomain creates a new domain in libvirt
//
// Supports rollback operation, safe to use in anywhere, anytime
func (s *ClientImpl) CreateDomain(ctx context.Context, domain Domain) error {
	conn, err := s.getConnect()
	if err != nil {
		return err
	}

	// Create new qcow2 image from base image
	if err = s.CreateImage(ctx, CreateImageParams{
		BaseImagePath:  domain.BaseImagePath(),
		CloneImagePath: domain.VMImagePath(),
		Size:           domain.Storage,
	}); err != nil {
		return fmt.Errorf("failed to clone image: %v", err)
	}

	domainXML, err := getXMLConfig(domain)
	if err != nil {
		return fmt.Errorf("failed to generate domain XML: %v", err)
	}

	xmlData, err := domainXML.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal domain XML: %v", err)
	}

	_, err = conn.DomainDefineXML(xmlData)
	if err != nil {
		return fmt.Errorf("failed to define domain: %v", err)
	}

	return nil
}

type UpdateDomainParams struct {
	Name    *string
	Cpu     *uint
	Ram     *uint
	Storage *uint
}

func (s *ClientImpl) UpdateDomain(ctx context.Context, domainID string, params UpdateDomainParams) error {
	conn, err := s.getConnect()
	if err != nil {
		return err
	}

	libDomain, err := s.getDomain(domainID)
	if err != nil {
		return err
	}

	domain, err := FromLibvirtToDomain(*libDomain)
	if err != nil {
		return fmt.Errorf("failed to convert domain to model: %v", err)
	}

	// Start updating domain
	if params.Name != nil {
		domain.Name = *params.Name
	}
	if params.Cpu != nil {
		domain.Cpu.Value = *params.Cpu
	}
	if params.Ram != nil {
		domain.Memory.Value = *params.Ram
	}
	if params.Storage != nil {
		domain.Storage = *params.Storage
	}

	// Generate new XML config
	domainXML, err := getXMLConfig(domain)
	if err != nil {
		return fmt.Errorf("failed to generate domain XML: %v", err)
	}

	xmlData, err := domainXML.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal domain XML: %v", err)
	}

	_, err = conn.DomainDefineXML(xmlData)
	if err != nil {
		return fmt.Errorf("failed to define domain: %v", err)
	}

	// TODO: real update domain, not redefine and start it again
	// If the domain is running, we need to start it again
	if val, err := libDomain.IsActive(); err != nil && val {
		if err := libDomain.Destroy(); err != nil {
			return fmt.Errorf("failed to destroy domain: %v", err)
		}
	}

	return nil
}

// DeleteDomain removes the domain from libvirt
//
// Does not support rollback operation so it should done last
func (s *ClientImpl) DeleteDomain(ctx context.Context, domainID string) error {
	libDomain, err := s.getDomain(domainID)
	if err != nil {
		return err
	}

	domain, err := FromLibvirtToDomain(*libDomain)
	if err != nil {
		return fmt.Errorf("failed to convert domain to model: %v", err)
	}

	isActive, err := libDomain.IsActive()
	if err != nil {
		return fmt.Errorf("failed to check if domain is active: %v", err)
	}

	if isActive {
		if err := libDomain.Destroy(); err != nil {
			return fmt.Errorf("failed to destroy domain: %v", err)
		}
	}

	if err = libDomain.Undefine(); err != nil {
		return fmt.Errorf("failed to undefine domain: %v", err)
	}

	// remove vm and cloudinit (always after domain is stopped and should not return error)
	//! These removal operations cannot be rolled back, so it should be done last
	if err := os.Remove(domain.VMImagePath()); err != nil {
		logger.Log.Error("failed to remove vm image", zap.String("path", domain.VMImagePath()), zap.Error(err))
	}
	if err := os.Remove(domain.CloudinitPath()); err != nil {
		logger.Log.Error("failed to remove cloudinit", zap.String("path", domain.CloudinitPath()), zap.Error(err))
	}

	return nil
}

func (s *ClientImpl) StartDomain(ctx context.Context, domainID string) error {
	domain, err := s.getDomain(domainID)
	if err != nil {
		return err
	}

	return domain.Create()
}

func (s *ClientImpl) StopDomain(ctx context.Context, domainID string) error {
	domain, err := s.getDomain(domainID)
	if err != nil {
		return err
	}

	return domain.Shutdown()
}
