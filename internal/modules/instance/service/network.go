package instancesvc

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	nginx "github.com/wagecloud/wagecloud-server/internal/client/nginx"
	instancemodel "github.com/wagecloud/wagecloud-server/internal/modules/instance/model"
	instancestorage "github.com/wagecloud/wagecloud-server/internal/modules/instance/storage"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
)

type MapPortNginxParams struct {
	VMIP         string
	ExternalPort int32
	InternalPort int32
	Type         string // "stream" or "http"
}

func (s *ServiceImpl) MapPortNginx(ctx context.Context, params MapPortNginxParams) error {
	// err := AddOrUpdateServerBlock(filepath.Join(os.Getenv("HOME"), "my-nginx/users.d/stream/test.conf"), "192.168.122.235", 22, 2345, stream)
	var pathName string

	if params.Type == "stream" {
		pathName = filepath.Join(os.Getenv("HOME"), "my-nginx/users.d/stream/test.conf")
	} else if params.Type == "http" {
		pathName = filepath.Join(os.Getenv("HOME"), "my-nginx/users.d/http/test.conf")
	} else {
		return fmt.Errorf("unsupported protocol type: %s", params.Type)
	}

	err := nginx.AddOrUpdateServerBlock(nginx.AddOrUpdateServerBlockParams{
		PathName:     pathName,
		VMIP:         params.VMIP,
		InternalPort: int(params.InternalPort),
		HostPort:     int(params.ExternalPort),
		ProtocolType: nginx.ProtocolType(params.Type),
	})

	if err != nil {
		return fmt.Errorf("error adding or updating server block: %w", err)
	}

	err = nginx.Reloading()
	if err != nil {
		return fmt.Errorf("error reloading nginx: %w", err)
	}
	return nil
}

type UnmapPortNginxParams struct {
	ExternalPort int32
	ProtocolType string // "stream" or "http"
}

func (s *ServiceImpl) UnmapPortNginx(ctx context.Context, params UnmapPortNginxParams) error {
	externalPort := params.ExternalPort
	protocolType := params.ProtocolType

	var pathName string

	if protocolType == "stream" {
		pathName = filepath.Join(os.Getenv("HOME"), "my-nginx/users.d/stream/test.conf")
	} else if protocolType == "http" {
		pathName = filepath.Join(os.Getenv("HOME"), "my-nginx/users.d/http/test.conf")
	} else {
		return fmt.Errorf("unsupported protocol type: %s", protocolType)
	}

	hostPortStr := fmt.Sprintf("%d", externalPort)
	err := nginx.DeleteServerBlock(pathName, hostPortStr)

	if err != nil {
		return fmt.Errorf("error deleting server block: %w", err)
	}

	err = nginx.Reloading()
	if err != nil {
		return fmt.Errorf("error reloading nginx: %w", err)
	}
	return nil
}

type GetNetworkParams struct {
	ID         *int64
	InstanceID *string
}

func (s *ServiceImpl) GetNetwork(ctx context.Context, params GetNetworkParams) (instancemodel.Network, error) {
	return s.storage.GetNetwork(ctx, instancestorage.GetNetworkParams{
		ID:         params.ID,
		InstanceID: params.InstanceID,
	})
}

type ListNetworksParams struct {
	pagination.PaginationParams
	ID        *string
	PrivateIP *string
	PublicIP  *string
}

func (s *ServiceImpl) ListNetworks(ctx context.Context, params ListNetworksParams) (res pagination.PaginateResult[instancemodel.Network], err error) {
	repoParams := instancestorage.ListNetworksParams{
		PaginationParams: params.PaginationParams,
		InstanceID:       params.ID,
		PrivateIP:        params.PrivateIP,
		PublicIP:         params.PublicIP,
	}

	total, err := s.storage.CountNetworks(ctx, repoParams)
	if err != nil {
		return res, err
	}

	networks, err := s.storage.ListNetworks(ctx, repoParams)
	if err != nil {
		return res, err
	}

	return pagination.PaginateResult[instancemodel.Network]{
		Data:     networks,
		Limit:    params.Limit,
		Page:     params.Page,
		Total:    total,
		NextPage: params.NextPage(total),
	}, nil
}

type CreateNetworkParams struct {
	InstanceID string
	PrivateIP  string
	PublicIP   *string
}

func (s *ServiceImpl) CreateNetwork(ctx context.Context, params CreateNetworkParams) (instancemodel.Network, error) {
	return s.storage.CreateNetwork(ctx, instancemodel.Network{
		InstanceID: params.InstanceID,
		PrivateIP:  params.PrivateIP,
		PublicIP:   params.PublicIP,
	})
}

type UpdateNetworkParams struct {
	ID           *int64
	InstanceID   *string
	PrivateIP    *string
	PublicIP     *string
	NullPublicIP bool
}

func (s *ServiceImpl) UpdateNetwork(ctx context.Context, params UpdateNetworkParams) (instancemodel.Network, error) {
	return s.storage.UpdateNetwork(ctx, instancestorage.UpdateNetworkParams{
		ID:           params.ID,
		InstanceID:   params.InstanceID,
		PrivateIP:    params.PrivateIP,
		PublicIP:     params.PublicIP,
		NullPublicIP: params.NullPublicIP,
	})
}

type DeleteNetworkParams struct {
	ID int64
}

func (s *ServiceImpl) DeleteNetwork(ctx context.Context, params DeleteNetworkParams) error {
	return s.storage.DeleteNetwork(ctx, params.ID)
}
