package instanceconnect

import (
	"context"

	"connectrpc.com/connect"
	instancev1 "github.com/wagecloud/wagecloud-server/gen/pb/instance/v1"
	instancemodel "github.com/wagecloud/wagecloud-server/internal/modules/instance/model"
	instancesvc "github.com/wagecloud/wagecloud-server/internal/modules/instance/service"
	commonmodel "github.com/wagecloud/wagecloud-server/internal/shared/model"
	"github.com/wagecloud/wagecloud-server/internal/utils/slice"
)

// Get network by ID
func (t *ImplementedInstanceServiceHandler) GetNetwork(ctx context.Context, req *connect.Request[instancev1.GetNetworkRequest]) (*connect.Response[instancev1.GetNetworkResponse], error) {
	result, err := t.service.GetNetwork(ctx, instancesvc.GetNetworkParams{
		ID: req.Msg.Id,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&instancev1.GetNetworkResponse{
		Network: instancemodel.NetworkModelToProto(result),
	}), nil
}

// List networks
func (t *ImplementedInstanceServiceHandler) ListNetworks(ctx context.Context, req *connect.Request[instancev1.ListNetworksRequest]) (*connect.Response[instancev1.ListNetworksResponse], error) {
	result, err := t.service.ListNetworks(ctx, instancesvc.ListNetworksParams{
		PaginationParams: commonmodel.PaginationParamsProtoToModel(req.Msg.Pagination),
		ID:               req.Msg.Id,
		PrivateIP:        req.Msg.PrivateIp,
		CreatedAtFrom:    req.Msg.CreatedAtFrom,
		CreatedAtTo:      req.Msg.CreatedAtTo,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&instancev1.ListNetworksResponse{
		Networks:   slice.Map(result.Data, instancemodel.NetworkModelToProto),
		Pagination: commonmodel.PaginateResultModelToProto(result),
	}), nil
}

// Create network
func (t *ImplementedInstanceServiceHandler) CreateNetwork(ctx context.Context, req *connect.Request[instancev1.CreateNetworkRequest]) (*connect.Response[instancev1.CreateNetworkResponse], error) {
	result, err := t.service.CreateNetwork(ctx, instancesvc.CreateNetworkParams{
		ID:        req.Msg.Id,
		PrivateIP: req.Msg.PrivateIp,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&instancev1.CreateNetworkResponse{
		Network: instancemodel.NetworkModelToProto(result),
	}), nil
}

// Update network
func (t *ImplementedInstanceServiceHandler) UpdateNetwork(ctx context.Context, req *connect.Request[instancev1.UpdateNetworkRequest]) (*connect.Response[instancev1.UpdateNetworkResponse], error) {
	result, err := t.service.UpdateNetwork(ctx, instancesvc.UpdateNetworkParams{
		ID:        req.Msg.Id,
		PrivateIP: req.Msg.PrivateIp,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&instancev1.UpdateNetworkResponse{
		Network: instancemodel.NetworkModelToProto(result),
	}), nil
}

// Delete network
func (t *ImplementedInstanceServiceHandler) DeleteNetwork(ctx context.Context, req *connect.Request[instancev1.DeleteNetworkRequest]) (*connect.Response[instancev1.DeleteNetworkResponse], error) {
	err := t.service.DeleteNetwork(ctx, instancesvc.DeleteNetworkParams{
		ID: req.Msg.Id,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&instancev1.DeleteNetworkResponse{}), nil
}
