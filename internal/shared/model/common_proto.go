package commonmodel

import (
	commonv1 "github.com/wagecloud/wagecloud-server/gen/pb/common/v1"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
)

func PaginateResultModelToProto[T any](result pagination.PaginateResult[T]) *commonv1.PaginateResult {
	return &commonv1.PaginateResult{
		Page:       result.Page,
		Limit:      result.Limit,
		Total:      result.Total,
		NextPage:   result.NextPage,
		NextCursor: result.NextCursor,
	}
}

func PaginateResultProtoToModel[T any](result *commonv1.PaginateResult, data []T) pagination.PaginateResult[T] {
	return pagination.PaginateResult[T]{
		Data:       data,
		Page:       result.Page,
		Limit:      result.Limit,
		Total:      result.Total,
		NextPage:   result.NextPage,
		NextCursor: result.NextCursor,
	}
}

func PaginationParamsProtoToModel(params *commonv1.PaginationParams) pagination.PaginationParams {
	if params == nil {
		return pagination.PaginationParams{}
	}
	return pagination.PaginationParams{
		Page:  params.Page,
		Limit: params.Limit,
	}
}

func PaginationParamsModelToProto(params pagination.PaginationParams) *commonv1.PaginationParams {
	return &commonv1.PaginationParams{
		Page:  params.Page,
		Limit: params.Limit,
	}
}
