package paymentstorage

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/wagecloud/wagecloud-server/gen/sqlc"
	"github.com/wagecloud/wagecloud-server/internal/client/pgxpool"
	paymentmodel "github.com/wagecloud/wagecloud-server/internal/modules/payment/model"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
	pgxptr "github.com/wagecloud/wagecloud-server/internal/utils/pgx/ptr"
	"github.com/wagecloud/wagecloud-server/internal/utils/ptr"
)

type Storage struct {
	db   pgxpool.DBTX
	sqlc *sqlc.Queries
}

func NewStorage(db pgxpool.DBTX) *Storage {
	return &Storage{
		db:   db,
		sqlc: sqlc.New(db),
	}
}

func (s *Storage) GetPayment(ctx context.Context, id int64) (paymentmodel.PaymentBase, error) {
	payment, err := s.sqlc.GetPayment(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return paymentmodel.PaymentBase{}, nil
		}
		return paymentmodel.PaymentBase{}, err
	}

	return paymentmodel.PaymentBase{
		ID:          payment.ID,
		AccountID:   payment.AccountID,
		Method:      paymentmodel.PaymentMethod(payment.Method),
		Status:      paymentmodel.PaymentStatus(payment.Status),
		Total:       payment.Total,
		DateCreated: payment.DateCreated.Time.UnixMilli(),
	}, nil
}

type ListPaymentsParams struct {
	pagination.PaginationParams
	ID              *string
	AccountID       *int64
	Method          *paymentmodel.PaymentMethod
	Status          *paymentmodel.PaymentStatus
	DateCreatedFrom *int64
	DateCreatedTo   *int64
}

func (s *Storage) CountPayments(ctx context.Context, params ListPaymentsParams) (int64, error) {
	return s.sqlc.CountPayments(ctx, sqlc.CountPaymentsParams{
		ID:              *pgxptr.PtrToPgtype(&pgtype.Text{}, params.ID),
		AccountID:       *pgxptr.PtrToPgtype(&pgtype.Int8{}, params.AccountID),
		Method:          *pgxptr.PtrBrandedToPgType(&sqlc.NullPaymentMethod{}, params.Method),
		Status:          *pgxptr.PtrBrandedToPgType(&sqlc.NullPaymentStatus{}, params.Status),
		DateCreatedFrom: *pgxptr.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.DateCreatedFrom)),
		DateCreatedTo:   *pgxptr.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.DateCreatedTo)),
	})
}

func (s *Storage) ListPayments(ctx context.Context, params ListPaymentsParams) ([]paymentmodel.PaymentBase, error) {
	payments, err := s.sqlc.ListPayments(ctx, sqlc.ListPaymentsParams{
		ID:              *pgxptr.PtrToPgtype(&pgtype.Text{}, params.ID),
		AccountID:       *pgxptr.PtrToPgtype(&pgtype.Int8{}, params.AccountID),
		Method:          *pgxptr.PtrBrandedToPgType(&sqlc.NullPaymentMethod{}, params.Method),
		Status:          *pgxptr.PtrBrandedToPgType(&sqlc.NullPaymentStatus{}, params.Status),
		DateCreatedFrom: *pgxptr.PtrToPgtype(&pgtype.Timestamptz{}, params.DateCreatedFrom),
		DateCreatedTo:   *pgxptr.PtrToPgtype(&pgtype.Timestamptz{}, params.DateCreatedTo),
		Offset:          params.Offset(),
		Limit:           params.Limit,
	})
	if err != nil {
		return nil, err
	}

	result := make([]paymentmodel.PaymentBase, len(payments))
	for i, payment := range payments {
		result[i] = paymentmodel.PaymentBase{
			ID:          payment.ID,
			AccountID:   payment.AccountID,
			Method:      paymentmodel.PaymentMethod(payment.Method),
			Status:      paymentmodel.PaymentStatus(payment.Status),
			Total:       payment.Total,
			DateCreated: payment.DateCreated.Time.UnixMilli(),
		}
	}

	return result, nil
}

func (s *Storage) CreatePayment(ctx context.Context, payment paymentmodel.PaymentBase) (paymentmodel.PaymentBase, error) {
	result, err := s.sqlc.CreatePayment(ctx, sqlc.CreatePaymentParams{
		AccountID: payment.AccountID,
		Method:    sqlc.PaymentMethod(payment.Method),
		Status:    sqlc.PaymentStatus(payment.Status),
		Total:     payment.Total,
	})
	if err != nil {
		return paymentmodel.PaymentBase{}, err
	}

	return paymentmodel.PaymentBase{
		ID:          result.ID,
		AccountID:   result.AccountID,
		Method:      paymentmodel.PaymentMethod(result.Method),
		Status:      paymentmodel.PaymentStatus(result.Status),
		Total:       result.Total,
		DateCreated: result.DateCreated.Time.UnixMilli(),
	}, nil
}

type UpdatePaymentParams struct {
	ID     int64
	Method *paymentmodel.PaymentMethod
	Status *paymentmodel.PaymentStatus
	Total  *int64
}

func (s *Storage) UpdatePayment(ctx context.Context, params UpdatePaymentParams) (paymentmodel.PaymentBase, error) {
	row, err := s.sqlc.UpdatePayment(ctx, sqlc.UpdatePaymentParams{
		ID:     params.ID,
		Method: *pgxptr.PtrBrandedToPgType(&sqlc.NullPaymentMethod{}, params.Method),
		Status: *pgxptr.PtrBrandedToPgType(&sqlc.NullPaymentStatus{}, params.Status),
		Total:  *pgxptr.PtrToPgtype(&pgtype.Int8{}, params.Total),
	})
	if err != nil {
		return paymentmodel.PaymentBase{}, err
	}

	return paymentmodel.PaymentBase{
		ID:          row.ID,
		AccountID:   row.AccountID,
		Method:      paymentmodel.PaymentMethod(row.Method),
		Status:      paymentmodel.PaymentStatus(row.Status),
		Total:       row.Total,
		DateCreated: row.DateCreated.Time.UnixMilli(),
	}, nil
}

func (s *Storage) DeletePayment(ctx context.Context, id int64) error {
	return s.sqlc.DeletePayment(ctx, id)
}

func (s *Storage) CreatePaymentItem(ctx context.Context, item paymentmodel.PaymentItem) (paymentmodel.PaymentItem, error) {
	row, err := s.sqlc.CreatePaymentItem(ctx, sqlc.CreatePaymentItemParams{
		PaymentID: item.PaymentID,
		Name:      item.Name,
		Price:     item.Price,
	})
	if err != nil {
		return paymentmodel.PaymentItem{}, err
	}

	return paymentmodel.PaymentItem{
		ID:        row.ID,
		PaymentID: row.PaymentID,
		Name:      row.Name,
		Price:     row.Price,
	}, nil
}

func (s *Storage) CreatePaymentVNPAY(ctx context.Context, vnpay paymentmodel.PaymentVNPAY) (paymentmodel.PaymentVNPAY, error) {
	row, err := s.sqlc.CreatePaymentVnpay(ctx, sqlc.CreatePaymentVnpayParams{
		ID:                 vnpay.ID,
		VnpTxnRef:          vnpay.VnpTxnRef,
		VnpOrderInfo:       vnpay.VnpOrderInfo,
		VnpTransactionNo:   vnpay.VnpTransactionNo,
		VnpTransactionDate: vnpay.VnpTransactionDate,
		VnpCreateDate:      vnpay.VnpCreateDate,
		VnpIpAddr:          vnpay.VnpIpAddr,
	})
	if err != nil {
		return paymentmodel.PaymentVNPAY{}, err
	}

	return paymentmodel.PaymentVNPAY{
		ID:                 row.ID,
		VnpTxnRef:          row.VnpTxnRef,
		VnpOrderInfo:       row.VnpOrderInfo,
		VnpTransactionNo:   row.VnpTransactionNo,
		VnpTransactionDate: row.VnpTransactionDate,
		VnpCreateDate:      row.VnpCreateDate,
		VnpIpAddr:          row.VnpIpAddr,
	}, nil
}
