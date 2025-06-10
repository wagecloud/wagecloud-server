package ossvc

import (
	"context"

	"connectrpc.com/connect"
	commonv1 "github.com/wagecloud/wagecloud-server/gen/pb/common/v1"
	osv1 "github.com/wagecloud/wagecloud-server/gen/pb/os/v1"
	osmodel "github.com/wagecloud/wagecloud-server/internal/modules/os/model"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
	"github.com/wagecloud/wagecloud-server/internal/utils/slice"
)

func (s *ServiceRpcImpl) GetArch(ctx context.Context, id string) (osmodel.Arch, error) {
	result, err := s.connect.GetArch(ctx, connect.NewRequest(&osv1.GetArchRequest{
		Id: id,
	}))
	if err != nil {
		return osmodel.Arch{}, err
	}

	return archProtoToModel(result.Msg.Arch), nil
}

func (s *ServiceRpcImpl) ListArchs(ctx context.Context, params ListArchsParams) (pagination.PaginateResult[osmodel.Arch], error) {
	result, err := s.connect.ListArchs(ctx, connect.NewRequest(&osv1.ListArchsRequest{
		Pagination: &commonv1.PaginationParams{
			Page:  params.Page,
			Limit: params.Limit,
		},
		Name:          params.Name,
		CreatedAtFrom: params.CreatedAtFrom,
		CreatedAtTo:   params.CreatedAtTo,
	}))
	if err != nil {
		return pagination.PaginateResult[osmodel.Arch]{}, err
	}

	return pagination.PaginateResult[osmodel.Arch]{
		Data:       slice.Map(result.Msg.Archs, archProtoToModel),
		Page:       result.Msg.Pagination.Page,
		Limit:      result.Msg.Pagination.Limit,
		Total:      result.Msg.Pagination.Total,
		NextPage:   result.Msg.Pagination.NextPage,
		NextCursor: result.Msg.Pagination.NextCursor,
	}, nil
}

func (s *ServiceRpcImpl) CreateArch(ctx context.Context, params CreateArchParams) (osmodel.Arch, error) {
	result, err := s.connect.CreateArch(ctx, connect.NewRequest(&osv1.CreateArchRequest{
		Id:   params.ID,
		Name: params.Name,
	}))
	if err != nil {
		return osmodel.Arch{}, err
	}

	return archProtoToModel(result.Msg.Arch), nil
}

func (s *ServiceRpcImpl) UpdateArch(ctx context.Context, params UpdateArchParams) (osmodel.Arch, error) {
	result, err := s.connect.UpdateArch(ctx, connect.NewRequest(&osv1.UpdateArchRequest{
		Id:    params.ID,
		NewId: params.NewID,
		Name:  params.Name,
	}))
	if err != nil {
		return osmodel.Arch{}, err
	}

	return archProtoToModel(result.Msg.Arch), nil
}

func (s *ServiceRpcImpl) DeleteArch(ctx context.Context, id string) error {
	_, err := s.connect.DeleteArch(ctx, connect.NewRequest(&osv1.DeleteArchRequest{
		Id: id,
	}))
	return err
}

func archProtoToModel(proto *osv1.Arch) osmodel.Arch {
	return osmodel.Arch{
		ID:        proto.Id,
		Name:      proto.Name,
		CreatedAt: proto.CreatedAt,
	}
}
