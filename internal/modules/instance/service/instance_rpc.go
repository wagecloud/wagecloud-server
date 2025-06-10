package instancesvc

import (
	"context"

	"connectrpc.com/connect"
	instancev1 "github.com/wagecloud/wagecloud-server/gen/pb/instance/v1"
	"github.com/wagecloud/wagecloud-server/gen/pb/instance/v1/instancev1connect"
	accountmodel "github.com/wagecloud/wagecloud-server/internal/modules/account/model"
	instancemodel "github.com/wagecloud/wagecloud-server/internal/modules/instance/model"
	commonmodel "github.com/wagecloud/wagecloud-server/internal/shared/model"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
	"github.com/wagecloud/wagecloud-server/internal/utils/slice"
)

type ServiceRpcImpl struct {
	connect instancev1connect.InstanceServiceClient
}

func NewServiceRpc(connect instancev1connect.InstanceServiceClient) Service {
	return &ServiceRpcImpl{
		connect: connect,
	}
}

func (s *ServiceRpcImpl) GetInstance(ctx context.Context, params GetInstanceParams) (instancemodel.Instance, error) {
	result, err := s.connect.GetInstance(ctx, connect.NewRequest(&instancev1.GetInstanceRequest{
		Id: params.ID,
	}))
	if err != nil {
		return instancemodel.Instance{}, err
	}

	return instancemodel.InstanceProtoToModel(result.Msg.Instance), nil
}

func (s *ServiceRpcImpl) ListInstances(ctx context.Context, params ListInstancesParams) (pagination.PaginateResult[instancemodel.Instance], error) {
	result, err := s.connect.ListInstances(ctx, connect.NewRequest(&instancev1.ListInstancesRequest{
		Pagination:    commonmodel.PaginationParamsModelToProto(params.PaginationParams),
		Account:       accountmodel.AuthenticatedAccountModelToProto(params.Account),
		OsId:          params.OsID,
		ArchId:        params.ArchID,
		Name:          params.Name,
		CpuFrom:       params.CpuFrom,
		CpuTo:         params.CpuTo,
		RamFrom:       params.RamFrom,
		RamTo:         params.RamTo,
		StorageFrom:   params.StorageFrom,
		StorageTo:     params.StorageTo,
		CreatedAtFrom: params.CreatedAtFrom,
		CreatedAtTo:   params.CreatedAtTo,
	}))
	if err != nil {
		return pagination.PaginateResult[instancemodel.Instance]{}, err
	}

	return pagination.PaginateResult[instancemodel.Instance]{
		Data:       slice.Map(result.Msg.Instances, instancemodel.InstanceProtoToModel),
		Page:       result.Msg.Pagination.Page,
		Limit:      result.Msg.Pagination.Limit,
		Total:      result.Msg.Pagination.Total,
		NextPage:   result.Msg.Pagination.NextPage,
		NextCursor: result.Msg.Pagination.NextCursor,
	}, nil
}

func (s *ServiceRpcImpl) CreateInstance(ctx context.Context, params CreateInstanceParams) (instancemodel.Instance, error) {
	result, err := s.connect.CreateInstance(ctx, connect.NewRequest(&instancev1.CreateInstanceRequest{
		Account:           accountmodel.AuthenticatedAccountModelToProto(params.Account),
		Name:              params.Name,
		SshAuthorizedKeys: params.SSHAuthorizedKeys,
		Password:          params.Password,
		LocalHostname:     params.LocalHostname,
		OsId:              params.OsID,
		ArchId:            params.ArchID,
		Memory:            params.Memory,
		Cpu:               params.Cpu,
		Storage:           params.Storage,
	}))
	if err != nil {
		return instancemodel.Instance{}, err
	}

	return instancemodel.InstanceProtoToModel(result.Msg.Instance), nil
}

func (s *ServiceRpcImpl) UpdateInstance(ctx context.Context, params UpdateInstanceParams) (instancemodel.Instance, error) {
	result, err := s.connect.UpdateInstance(ctx, connect.NewRequest(&instancev1.UpdateInstanceRequest{
		Id:        params.ID,
		Name:      params.Name,
		NetworkId: params.NetworkID,
		OsId:      params.OsID,
		ArchId:    params.ArchID,
		Cpu:       params.Cpu,
		Ram:       params.Ram,
		Storage:   params.Storage,
	}))
	if err != nil {
		return instancemodel.Instance{}, err
	}

	return instancemodel.InstanceProtoToModel(result.Msg.Instance), nil
}

func (s *ServiceRpcImpl) DeleteInstance(ctx context.Context, params DeleteInstanceParams) error {
	_, err := s.connect.DeleteInstance(ctx, connect.NewRequest(&instancev1.DeleteInstanceRequest{
		Id: params.ID,
	}))
	return err
}

func (s *ServiceRpcImpl) StartInstance(ctx context.Context, params StartInstanceParams) error {
	_, err := s.connect.StartInstance(ctx, connect.NewRequest(&instancev1.StartInstanceRequest{
		Id: params.ID,
	}))
	return err
}

func (s *ServiceRpcImpl) StopInstance(ctx context.Context, params StopInstanceParams) error {
	_, err := s.connect.StopInstance(ctx, connect.NewRequest(&instancev1.StopInstanceRequest{
		Id: params.ID,
	}))
	return err
}

func (s *ServiceRpcImpl) DeleteNetwork(ctx context.Context, params DeleteNetworkParams) error {
	_, err := s.connect.DeleteNetwork(ctx, connect.NewRequest(&instancev1.DeleteNetworkRequest{
		Id: params.ID,
	}))
	return err
}
