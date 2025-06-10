package ossvc

import (
	"context"

	"connectrpc.com/connect"
	commonv1 "github.com/wagecloud/wagecloud-server/gen/pb/common/v1"
	osv1 "github.com/wagecloud/wagecloud-server/gen/pb/os/v1"
	"github.com/wagecloud/wagecloud-server/gen/pb/os/v1/osv1connect"
	osmodel "github.com/wagecloud/wagecloud-server/internal/modules/os/model"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
	"github.com/wagecloud/wagecloud-server/internal/utils/slice"
)

type ServiceRpcImpl struct {
	connect osv1connect.OSServiceClient
}

func NewServiceRpc(connect osv1connect.OSServiceClient) Service {
	return &ServiceRpcImpl{
		connect: connect,
	}
}

func (s *ServiceRpcImpl) GetOS(ctx context.Context, params GetOSParams) (osmodel.OS, error) {
	result, err := s.connect.GetOS(ctx, connect.NewRequest(&osv1.GetOSRequest{
		Id: params.ID,
	}))
	if err != nil {
		return osmodel.OS{}, err
	}

	return osmodel.OSProtoToModel(result.Msg.Os), nil
}

func (s *ServiceRpcImpl) ListOSs(ctx context.Context, params ListOSsParams) (pagination.PaginateResult[osmodel.OS], error) {
	result, err := s.connect.ListOSs(ctx, connect.NewRequest(&osv1.ListOSsRequest{
		Pagination: &commonv1.PaginationParams{
			Page:  params.Page,
			Limit: params.Limit,
		},
		Name:          params.Name,
		CreatedAtFrom: params.CreatedAtFrom,
		CreatedAtTo:   params.CreatedAtTo,
	}))
	if err != nil {
		return pagination.PaginateResult[osmodel.OS]{}, err
	}

	return pagination.PaginateResult[osmodel.OS]{
		Data:       slice.Map(result.Msg.Oss, osmodel.OSProtoToModel),
		Page:       result.Msg.Pagination.Page,
		Limit:      result.Msg.Pagination.Limit,
		Total:      result.Msg.Pagination.Total,
		NextPage:   result.Msg.Pagination.NextPage,
		NextCursor: result.Msg.Pagination.NextCursor,
	}, nil
}

func (s *ServiceRpcImpl) CreateOS(ctx context.Context, params CreateOSParams) (osmodel.OS, error) {
	result, err := s.connect.CreateOS(ctx, connect.NewRequest(&osv1.CreateOSRequest{
		Id:   params.ID,
		Name: params.Name,
	}))
	if err != nil {
		return osmodel.OS{}, err
	}

	return osmodel.OSProtoToModel(result.Msg.Os), nil
}

func (s *ServiceRpcImpl) UpdateOS(ctx context.Context, params UpdateOSParams) (osmodel.OS, error) {
	result, err := s.connect.UpdateOS(ctx, connect.NewRequest(&osv1.UpdateOSRequest{
		Id:    params.ID,
		NewId: params.NewID,
		Name:  params.Name,
	}))
	if err != nil {
		return osmodel.OS{}, err
	}

	return osmodel.OSProtoToModel(result.Msg.Os), nil
}

func (s *ServiceRpcImpl) DeleteOS(ctx context.Context, params DeleteOSParams) error {
	_, err := s.connect.DeleteOS(ctx, connect.NewRequest(&osv1.DeleteOSRequest{
		Id: params.ID,
	}))
	return err
}
