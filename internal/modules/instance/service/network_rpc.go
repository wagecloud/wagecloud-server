package instancesvc

import (
	"context"

	"connectrpc.com/connect"
	commonv1 "github.com/wagecloud/wagecloud-server/gen/pb/common/v1"
	instancev1 "github.com/wagecloud/wagecloud-server/gen/pb/instance/v1"
	instancemodel "github.com/wagecloud/wagecloud-server/internal/modules/instance/model"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
	"github.com/wagecloud/wagecloud-server/internal/utils/slice"
)

func (s *ServiceRpcImpl) GetNetwork(ctx context.Context, params GetNetworkParams) (instancemodel.Network, error) {
	result, err := s.connect.GetNetwork(ctx, connect.NewRequest(&instancev1.GetNetworkRequest{
		Id: params.ID,
	}))
	if err != nil {
		return instancemodel.Network{}, err
	}

	return networkProtoToModel(result.Msg.Network), nil
}

func (s *ServiceRpcImpl) ListNetworks(ctx context.Context, params ListNetworksParams) (pagination.PaginateResult[instancemodel.Network], error) {
	result, err := s.connect.ListNetworks(ctx, connect.NewRequest(&instancev1.ListNetworksRequest{
		Pagination: &commonv1.PaginationParams{
			Page:  params.Page,
			Limit: params.Limit,
		},
		Id:            params.ID,
		PrivateIp:     params.PrivateIP,
		CreatedAtFrom: params.CreatedAtFrom,
		CreatedAtTo:   params.CreatedAtTo,
	}))
	if err != nil {
		return pagination.PaginateResult[instancemodel.Network]{}, err
	}

	return pagination.PaginateResult[instancemodel.Network]{
		Data:       slice.Map(result.Msg.Networks, networkProtoToModel),
		Page:       result.Msg.Pagination.Page,
		Limit:      result.Msg.Pagination.Limit,
		Total:      result.Msg.Pagination.Total,
		NextPage:   result.Msg.Pagination.NextPage,
		NextCursor: result.Msg.Pagination.NextCursor,
	}, nil
}

func (s *ServiceRpcImpl) CreateNetwork(ctx context.Context, params CreateNetworkParams) (instancemodel.Network, error) {
	result, err := s.connect.CreateNetwork(ctx, connect.NewRequest(&instancev1.CreateNetworkRequest{
		Id:        params.ID,
		PrivateIp: params.PrivateIP,
	}))
	if err != nil {
		return instancemodel.Network{}, err
	}

	return networkProtoToModel(result.Msg.Network), nil
}

func (s *ServiceRpcImpl) UpdateNetwork(ctx context.Context, params UpdateNetworkParams) (instancemodel.Network, error) {
	result, err := s.connect.UpdateNetwork(ctx, connect.NewRequest(&instancev1.UpdateNetworkRequest{
		Id:        params.ID,
		NewId:     params.NewID,
		PrivateIp: params.PrivateIP,
	}))
	if err != nil {
		return instancemodel.Network{}, err
	}

	return networkProtoToModel(result.Msg.Network), nil
}

func networkProtoToModel(proto *instancev1.Network) instancemodel.Network {
	return instancemodel.Network{
		ID:        proto.Id,
		PrivateIP: proto.PrivateIp,
		CreatedAt: proto.CreatedAt,
	}
}
