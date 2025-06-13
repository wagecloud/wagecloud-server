package osconnect

import (
	"context"

	"connectrpc.com/connect"
	commonv1 "github.com/wagecloud/wagecloud-server/gen/pb/common/v1"
	osv1 "github.com/wagecloud/wagecloud-server/gen/pb/os/v1"
	"github.com/wagecloud/wagecloud-server/gen/pb/os/v1/osv1connect"
	osmodel "github.com/wagecloud/wagecloud-server/internal/modules/os/model"
	ossvc "github.com/wagecloud/wagecloud-server/internal/modules/os/service"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
	"github.com/wagecloud/wagecloud-server/internal/utils/slice"
)

type ImplementedOSServiceHandler struct {
	osv1connect.UnimplementedOSServiceHandler
	service ossvc.Service
}

func NewImplementedOSServiceHandler(service ossvc.Service) osv1connect.OSServiceHandler {
	return &ImplementedOSServiceHandler{
		service: service,
	}
}

func (t *ImplementedOSServiceHandler) GetOS(ctx context.Context, req *connect.Request[osv1.GetOSRequest]) (*connect.Response[osv1.GetOSResponse], error) {
	result, err := t.service.GetOS(ctx, ossvc.GetOSParams{
		ID: req.Msg.Id,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&osv1.GetOSResponse{
		Os: osmodel.OSModelToProto(result),
	}), nil
}

func (t *ImplementedOSServiceHandler) ListOSs(ctx context.Context, req *connect.Request[osv1.ListOSsRequest]) (*connect.Response[osv1.ListOSsResponse], error) {
	result, err := t.service.ListOSs(ctx, ossvc.ListOSsParams{
		PaginationParams: pagination.PaginationParams{
			Page:  req.Msg.Pagination.Page,
			Limit: req.Msg.Pagination.Limit,
		},
		Name:          req.Msg.Name,
		CreatedAtFrom: req.Msg.CreatedAtFrom,
		CreatedAtTo:   req.Msg.CreatedAtTo,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&osv1.ListOSsResponse{
		Oss: slice.Map(result.Data, osmodel.OSModelToProto),
		Pagination: &commonv1.PaginateResult{
			Page:       result.Page,
			Limit:      result.Limit,
			Total:      result.Total,
			NextPage:   result.NextPage,
			NextCursor: result.NextCursor,
		},
	}), nil
}

func (t *ImplementedOSServiceHandler) CreateOS(ctx context.Context, req *connect.Request[osv1.CreateOSRequest]) (*connect.Response[osv1.CreateOSResponse], error) {
	result, err := t.service.CreateOS(ctx, ossvc.CreateOSParams{
		ID:   req.Msg.Id,
		Name: req.Msg.Name,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&osv1.CreateOSResponse{
		Os: osmodel.OSModelToProto(result),
	}), nil
}

func (t *ImplementedOSServiceHandler) UpdateOS(ctx context.Context, req *connect.Request[osv1.UpdateOSRequest]) (*connect.Response[osv1.UpdateOSResponse], error) {
	result, err := t.service.UpdateOS(ctx, ossvc.UpdateOSParams{
		ID:    req.Msg.Id,
		NewID: req.Msg.NewId,
		Name:  req.Msg.Name,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&osv1.UpdateOSResponse{
		Os: osmodel.OSModelToProto(result),
	}), nil
}

func (t *ImplementedOSServiceHandler) DeleteOS(ctx context.Context, req *connect.Request[osv1.DeleteOSRequest]) (*connect.Response[osv1.DeleteOSResponse], error) {
	err := t.service.DeleteOS(ctx, ossvc.DeleteOSParams{
		ID: req.Msg.Id,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&osv1.DeleteOSResponse{}), nil
}

func (t *ImplementedOSServiceHandler) GetArch(ctx context.Context, req *connect.Request[osv1.GetArchRequest]) (*connect.Response[osv1.GetArchResponse], error) {
	result, err := t.service.GetArch(ctx, req.Msg.Id)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&osv1.GetArchResponse{
		Arch: osmodel.ArchModelToProto(result),
	}), nil
}

func (t *ImplementedOSServiceHandler) ListArchs(ctx context.Context, req *connect.Request[osv1.ListArchsRequest]) (*connect.Response[osv1.ListArchsResponse], error) {
	result, err := t.service.ListArchs(ctx, ossvc.ListArchsParams{
		PaginationParams: pagination.PaginationParams{
			Page:  req.Msg.Pagination.Page,
			Limit: req.Msg.Pagination.Limit,
		},
		Name:          req.Msg.Name,
		CreatedAtFrom: req.Msg.CreatedAtFrom,
		CreatedAtTo:   req.Msg.CreatedAtTo,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&osv1.ListArchsResponse{
		Archs: slice.Map(result.Data, osmodel.ArchModelToProto),
		Pagination: &commonv1.PaginateResult{
			Page:  result.Page,
			Limit: result.Limit,
			Total: result.Total,
		},
	}), nil
}

func (t *ImplementedOSServiceHandler) CreateArch(ctx context.Context, req *connect.Request[osv1.CreateArchRequest]) (*connect.Response[osv1.CreateArchResponse], error) {
	result, err := t.service.CreateArch(ctx, ossvc.CreateArchParams{
		ID:   req.Msg.Id,
		Name: req.Msg.Name,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&osv1.CreateArchResponse{
		Arch: osmodel.ArchModelToProto(result),
	}), nil
}

func (t *ImplementedOSServiceHandler) UpdateArch(ctx context.Context, req *connect.Request[osv1.UpdateArchRequest]) (*connect.Response[osv1.UpdateArchResponse], error) {
	result, err := t.service.UpdateArch(ctx, ossvc.UpdateArchParams{
		ID:    req.Msg.Id,
		NewID: req.Msg.NewId,
		Name:  req.Msg.Name,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&osv1.UpdateArchResponse{
		Arch: osmodel.ArchModelToProto(result),
	}), nil
}

func (t *ImplementedOSServiceHandler) DeleteArch(ctx context.Context, req *connect.Request[osv1.DeleteArchRequest]) (*connect.Response[osv1.DeleteArchResponse], error) {
	err := t.service.DeleteArch(ctx, req.Msg.Id)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&osv1.DeleteArchResponse{}), nil
}
