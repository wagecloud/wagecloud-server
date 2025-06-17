package instancestorage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/wagecloud/wagecloud-server/gen/sqlc"
	instancemodel "github.com/wagecloud/wagecloud-server/internal/modules/instance/model"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
	pgxptr "github.com/wagecloud/wagecloud-server/internal/utils/pgx/ptr"
)

func (r *Storage) GetInstanceLog(ctx context.Context, id int64) (instancemodel.InstanceLog, error) {
	row, err := r.sqlc.GetInstanceLog(ctx, id)
	if err != nil {
		return instancemodel.InstanceLog{}, err
	}

	return instancemodel.InstanceLog{
		ID:          row.ID,
		InstanceID:  row.InstanceID,
		Type:        instancemodel.LogType(row.Type),
		Title:       row.Title,
		Description: pgxptr.PgtypeToPtr[string](row.Description),
		CreatedAt:   row.CreatedAt.Time,
	}, nil
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

func (r *Storage) CountInstanceLogs(ctx context.Context, params ListInstanceLogsParams) (int64, error) {
	return r.sqlc.CountInstanceLogs(ctx, sqlc.CountInstanceLogsParams{
		InstanceID:    *pgxptr.PtrToPgtype(&pgtype.Text{}, params.InstanceID),
		Type:          *pgxptr.PtrBrandedToPgType(&sqlc.NullInstanceLogType{}, params.Type),
		Title:         *pgxptr.PtrToPgtype(&pgtype.Text{}, params.Title),
		Description:   *pgxptr.PtrToPgtype(&pgtype.Text{}, params.Description),
		CreatedAtFrom: *pgxptr.PtrToPgtype(&pgtype.Timestamptz{}, params.CreatedAtFrom),
		CreatedAtTo:   *pgxptr.PtrToPgtype(&pgtype.Timestamptz{}, params.CreatedAtTo),
	})
}

func (r *Storage) ListInstanceLogs(ctx context.Context, params ListInstanceLogsParams) ([]instancemodel.InstanceLog, error) {
	instanceLogs, err := r.sqlc.ListInstanceLogs(ctx, sqlc.ListInstanceLogsParams{
		Offset:        params.Offset(),
		Limit:         params.Limit,
		InstanceID:    *pgxptr.PtrToPgtype(&pgtype.Text{}, params.InstanceID),
		Type:          *pgxptr.PtrBrandedToPgType(&sqlc.NullInstanceLogType{}, params.Type),
		Title:         *pgxptr.PtrToPgtype(&pgtype.Text{}, params.Title),
		Description:   *pgxptr.PtrToPgtype(&pgtype.Text{}, params.Description),
		CreatedAtFrom: *pgxptr.PtrToPgtype(&pgtype.Timestamptz{}, params.CreatedAtFrom),
		CreatedAtTo:   *pgxptr.PtrToPgtype(&pgtype.Timestamptz{}, params.CreatedAtTo),
	})
	if err != nil {
		return nil, err
	}

	var result []instancemodel.InstanceLog
	for _, instanceLog := range instanceLogs {
		result = append(result, instancemodel.InstanceLog{
			ID:          instanceLog.ID,
			InstanceID:  instanceLog.InstanceID,
			Type:        instancemodel.LogType(instanceLog.Type),
			Title:       instanceLog.Title,
			Description: pgxptr.PgtypeToPtr[string](instanceLog.Description),
			CreatedAt:   instanceLog.CreatedAt.Time,
		})
	}

	return result, nil
}

func (r *Storage) CreateInstanceLog(ctx context.Context, instanceLog instancemodel.InstanceLog) (instancemodel.InstanceLog, error) {
	row, err := r.sqlc.CreateInstanceLog(ctx, sqlc.CreateInstanceLogParams{
		InstanceID:  instanceLog.InstanceID,
		Type:        sqlc.InstanceLogType(instanceLog.Type),
		Title:       instanceLog.Title,
		Description: *pgxptr.PtrToPgtype(&pgtype.Text{}, instanceLog.Description),
	})
	if err != nil {
		return instancemodel.InstanceLog{}, err
	}

	return instancemodel.InstanceLog{
		ID:          row.ID,
		InstanceID:  row.InstanceID,
		Type:        instancemodel.LogType(row.Type),
		Title:       row.Title,
		Description: pgxptr.PgtypeToPtr[string](row.Description),
		CreatedAt:   row.CreatedAt.Time,
	}, nil
}

type UpdateInstanceLogParams struct {
	ID              int64
	InstanceID      *string
	Type            *instancemodel.LogType
	Title           *string
	Description     *string
	NullDescription bool
}

func (r *Storage) UpdateInstanceLog(ctx context.Context, params UpdateInstanceLogParams) (instancemodel.InstanceLog, error) {
	row, err := r.sqlc.UpdateInstanceLog(ctx, sqlc.UpdateInstanceLogParams{
		ID:              params.ID,
		InstanceID:      *pgxptr.PtrToPgtype(&pgtype.Text{}, params.InstanceID),
		Type:            *pgxptr.PtrBrandedToPgType(&sqlc.NullInstanceLogType{}, params.Type),
		Title:           *pgxptr.PtrToPgtype(&pgtype.Text{}, params.Title),
		Description:     *pgxptr.PtrToPgtype(&pgtype.Text{}, params.Description),
		NullDescription: params.NullDescription,
	})
	if err != nil {
		return instancemodel.InstanceLog{}, err
	}

	return instancemodel.InstanceLog{
		ID:          row.ID,
		InstanceID:  row.InstanceID,
		Type:        instancemodel.LogType(row.Type),
		Title:       row.Title,
		Description: pgxptr.PgtypeToPtr[string](row.Description),
	}, nil
}

func (r *Storage) DeleteInstanceLog(ctx context.Context, id int64) error {
	return r.sqlc.DeleteInstanceLog(ctx, id)
}
