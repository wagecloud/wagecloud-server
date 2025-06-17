package instancesvc

import (
	"context"
	"time"

	instancemodel "github.com/wagecloud/wagecloud-server/internal/modules/instance/model"
	instancestorage "github.com/wagecloud/wagecloud-server/internal/modules/instance/storage"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
)

// TODO: add authenticatedAccount to every params!

func (s *ServiceImpl) GetInstanceLog(ctx context.Context, id int64) (instancemodel.InstanceLog, error) {
	return s.storage.GetInstanceLog(ctx, id)
}

type ListInstanceLogsParams struct {
	pagination.PaginationParams
	InstanceID    *string
	Type          *instancemodel.LogType
	Title         *string
	Description   *string
	CreatedAtFrom *time.Time
	CreatedAtTo   *time.Time
}

func (s *ServiceImpl) ListInstanceLogs(ctx context.Context, params ListInstanceLogsParams) (res pagination.PaginateResult[instancemodel.InstanceLog], err error) {
	// TODO: rename all repoParams into storageParams
	storageParams := instancestorage.ListInstanceLogsParams{
		PaginationParams: params.PaginationParams,
		InstanceID:       params.InstanceID,
		Type:             params.Type,
		Title:            params.Title,
		Description:      params.Description,
		CreatedAtFrom:    params.CreatedAtFrom,
		CreatedAtTo:      params.CreatedAtTo,
	}

	total, err := s.storage.CountInstanceLogs(ctx, storageParams)
	if err != nil {
		return res, err
	}

	instanceLogs, err := s.storage.ListInstanceLogs(ctx, storageParams)
	if err != nil {
		return res, err
	}

	return pagination.PaginateResult[instancemodel.InstanceLog]{
		Data:     instanceLogs,
		Limit:    params.Limit,
		Page:     params.Page,
		Total:    total,
		NextPage: params.NextPage(total),
	}, nil
}

type CreateInstanceLogParams struct {
	InstanceID  string
	Type        instancemodel.LogType
	Title       string
	Description *string
}

func (s *ServiceImpl) CreateInstanceLog(ctx context.Context, params CreateInstanceLogParams) (instancemodel.InstanceLog, error) {
	return s.storage.CreateInstanceLog(ctx, instancemodel.InstanceLog{
		InstanceID:  params.InstanceID,
		Type:        params.Type,
		Title:       params.Title,
		Description: params.Description,
	})
}

type UpdateInstanceLogParams struct {
	ID              int64
	Type            *instancemodel.LogType
	Title           *string
	Description     *string
	NullDescription bool
}

func (s *ServiceImpl) UpdateInstanceLog(ctx context.Context, params UpdateInstanceLogParams) (instancemodel.InstanceLog, error) {
	return s.storage.UpdateInstanceLog(ctx, instancestorage.UpdateInstanceLogParams{
		ID:              params.ID,
		Type:            params.Type,
		Title:           params.Title,
		Description:     params.Description,
		NullDescription: params.NullDescription,
	})
}

func (s *ServiceImpl) DeleteInstanceLog(ctx context.Context, id int64) error {
	return s.storage.DeleteInstanceLog(ctx, id)
}
