package libvirt

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"time"

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
	GetDomainMonitor(ctx context.Context, domainID string) (DomainMonitor, error)
	ListDomains(ctx context.Context, params ListDomainsParams) ([]Domain, error)
	CreateDomain(ctx context.Context, domain Domain) error
	UpdateDomain(ctx context.Context, domainID string, params UpdateDomainParams) error
	DeleteDomain(ctx context.Context, domainID string) error
	StartDomain(ctx context.Context, domainID string) error
	StopDomain(ctx context.Context, domainID string) error
	GetPrivateIP(ctx context.Context, domainID string) (string, error)

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

func (s *ClientImpl) GetPrivateIP(ctx context.Context, domainID string) (string, error) {
	domain, err := s.getDomain(domainID)
	if err != nil {
		return "", err
	}

	// Get the interface interfaces
	ifaces1, _ := domain.ListAllInterfaceAddresses(libvirt.DOMAIN_INTERFACE_ADDRESSES_SRC_ARP)
	ifaces2, _ := domain.ListAllInterfaceAddresses(libvirt.DOMAIN_INTERFACE_ADDRESSES_SRC_LEASE)
	ifaces3, _ := domain.ListAllInterfaceAddresses(libvirt.DOMAIN_INTERFACE_ADDRESSES_SRC_AGENT)
	ifaces := append(ifaces1, ifaces2...)
	ifaces = append(ifaces, ifaces3...)

	for _, iface := range ifaces {
		for _, addr := range iface.Addrs {
			if addr.Type == libvirt.IP_ADDR_TYPE_IPV4 {
				return addr.Addr, nil
			}
		}
	}

	return "", fmt.Errorf("no private IP found for domain %s", domainID)
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

func (s *ClientImpl) GetDomainMonitor(ctx context.Context, domainID string) (DomainMonitor, error) {
	conn, err := s.getConnect()
	if err != nil {
		return DomainMonitor{}, err
	}

	libDomain, err := conn.LookupDomainByUUIDString(domainID)
	if err != nil {
		return DomainMonitor{}, fmt.Errorf("failed to lookup domain by ID %s: %v", domainID, err)
	}
	defer libDomain.Free()

	// Get domain state
	state, _, err := libDomain.GetState()
	if err != nil {
		return DomainMonitor{}, fmt.Errorf("failed to get domain state: %v", err)
	}

	status := ToStatus(state)

	// Network stats
	// TODO: double check the list all interface address if it return more than two ifaces?
	// Get interface list once
	interfaces, _ := libDomain.ListAllInterfaceAddresses(libvirt.DOMAIN_INTERFACE_ADDRESSES_SRC_ARP)

	// First measurement
	var rx1, tx1 float64
	var rx2, tx2 float64

	var cpuUsagePercent float64

	if status == StatusPending || status == StatusRunning {
		for _, iface := range interfaces {
			stats, err := libDomain.InterfaceStats(iface.Name)
			if err == nil {
				rx1 += float64(stats.RxBytes) / (1024 * 1024)
				tx1 += float64(stats.TxBytes) / (1024 * 1024)
			}
		}

		// CPU stats
		// Sample 1
		startTime := time.Now()
		info, _ := libDomain.GetInfo()

		startStats, _ := libDomain.GetCPUStats(-1, 1, 0)
		startCpuTime := startStats[0].CpuTime

		// Wait some interval
		time.Sleep(1 * time.Second)

		// Sample 2
		endTime := time.Now()
		endStats, _ := libDomain.GetCPUStats(-1, 1, 0)
		endCpuTime := endStats[0].CpuTime

		// CPU usage calculation
		elapsed := endTime.Sub(startTime).Nanoseconds()
		used := endCpuTime - startCpuTime

		// Number of logical CPUs (can be retrieved from host or domain info)
		numCPUs := info.NrVirtCpu
		// or libDomain.GetInfo().NrVirtCpu

		cpuUsagePercent = float64(used) / float64(elapsed*int64(numCPUs)) * 100

		// Second measurement - use SAME interface list

		for _, iface := range interfaces {
			stats, err := libDomain.InterfaceStats(iface.Name)
			if err == nil {
				rx2 += float64(stats.RxBytes) / (1024 * 1024)
				tx2 += float64(stats.TxBytes) / (1024 * 1024)
			}
		}
	}

	// Memory stats
	memStats, err := libDomain.MemoryStats(13, 0)
	if err != nil {
		return DomainMonitor{}, fmt.Errorf("failed to get memory stats: %v", err)
	}

	// js, _ := json.Marshal(memStats)
	// fmt.Println("Memory Stats:", string(js))

	var availableMemory, unusedMemory uint64
	for _, stat := range memStats {
		if stat.Tag == int32(libvirt.DOMAIN_MEMORY_STAT_ACTUAL_BALLOON) {
			availableMemory = stat.Val // Total usable memory in KB
		}
		if stat.Tag == int32(libvirt.DOMAIN_MEMORY_STAT_USABLE) {
			unusedMemory = stat.Val // Actually used memory in KB
		}
	}

	// Calculate current memory consumption
	// var ramUsagePercent float64
	var ramUsageMB float64

	if availableMemory > 0 {
		ramUsageMB = float64(availableMemory-unusedMemory) / 1024
		// ramUsagePercent = (float64(availableMemory-unusedMemory) / float64(availableMemory)) * 100
	}

	// fmt.Printf("Memory Usage: %.2f MB (%.2f%%)\n", ramUsageMB, ramUsagePercent)

	// Block (disk) stats
	blockStats, err := libDomain.GetBlockInfo("vda", 0)
	var storageUsed float64
	if err == nil {
		storageUsed = float64(blockStats.Allocation) / (1024 * 1024) // Convert to MB
	}

	// Calculate throughput (MB/s)
	rxSpeed := math.Abs(rx2 - rx1) // MB/s over 1 second
	txSpeed := math.Abs(tx2 - tx1) // MB/s over 1 second
	return DomainMonitor{
		ID:           domainID,
		Status:       status,
		CPUUsage:     cpuUsagePercent,
		RAMUsage:     ramUsageMB,
		StorageUsage: storageUsed,
		NetworkIn:    rxSpeed,
		NetworkOut:   txSpeed,
	}, nil
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
