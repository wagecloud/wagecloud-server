package paymentstorage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/wagecloud/wagecloud-server/gen/sqlc"
	"github.com/wagecloud/wagecloud-server/internal/client/pgxpool"
	paymentmodel "github.com/wagecloud/wagecloud-server/internal/modules/payment/model"
	commonmodel "github.com/wagecloud/wagecloud-server/internal/shared/model"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
	pgxptr "github.com/wagecloud/wagecloud-server/internal/utils/pgx/ptr"
	"github.com/wagecloud/wagecloud-server/internal/utils/ptr"
)

type Storage struct {
	db   pgxpool.DBTX
	sqlc *sqlc.Queries
}

type TxStorage struct {
	*Storage
	tx pgx.Tx
}

func NewStorage(db pgxpool.DBTX) *Storage {
	return &Storage{
		db:   db,
		sqlc: sqlc.New(db),
	}
}

func (s *Storage) BeginTx(ctx context.Context) (*TxStorage, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return &TxStorage{
		Storage: NewStorage(tx),
		tx:      tx,
	}, nil
}

func (ts *TxStorage) Commit(ctx context.Context) error {
	return ts.tx.Commit(ctx)
}

func (ts *TxStorage) Rollback(ctx context.Context) error {
	return ts.tx.Rollback(ctx)
}

func (s *Storage) GetPayment(ctx context.Context, id int64) (paymentmodel.Payment, error) {
	payment, err := s.sqlc.GetPayment(ctx, id)
	if err != nil {
		return paymentmodel.Payment{}, err
	}

	return paymentmodel.Payment{
		ID:          payment.ID,
		AccountID:   payment.AccountID,
		Method:      paymentmodel.PaymentMethod(payment.Method),
		Status:      paymentmodel.PaymentStatus(payment.Status),
		Total:       commonmodel.Concurrency(payment.Total),
		DateCreated: payment.DateCreated.Time,
	}, nil
}

type ListPaymentsParams struct {
	pagination.PaginationParams
	AccountID       *int64
	Method          *paymentmodel.PaymentMethod
	Status          *paymentmodel.PaymentStatus
	DateCreatedFrom *int64
	DateCreatedTo   *int64
}

func (s *Storage) CountPayments(ctx context.Context, params ListPaymentsParams) (int64, error) {
	return s.sqlc.CountPayments(ctx, sqlc.CountPaymentsParams{
		AccountID:       *pgxptr.PtrToPgtype(&pgtype.Int8{}, params.AccountID),
		Method:          *pgxptr.PtrBrandedToPgType(&sqlc.NullPaymentMethod{}, params.Method),
		Status:          *pgxptr.PtrBrandedToPgType(&sqlc.NullPaymentStatus{}, params.Status),
		DateCreatedFrom: *pgxptr.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.DateCreatedFrom)),
		DateCreatedTo:   *pgxptr.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.DateCreatedTo)),
	})
}

func (s *Storage) ListPayments(ctx context.Context, params ListPaymentsParams) ([]paymentmodel.Payment, error) {
	payments, err := s.sqlc.ListPayments(ctx, sqlc.ListPaymentsParams{
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

	result := make([]paymentmodel.Payment, len(payments))
	for i, payment := range payments {
		result[i] = paymentmodel.Payment{
			ID:          payment.ID,
			AccountID:   payment.AccountID,
			Method:      paymentmodel.PaymentMethod(payment.Method),
			Status:      paymentmodel.PaymentStatus(payment.Status),
			Total:       commonmodel.Concurrency(payment.Total),
			DateCreated: payment.DateCreated.Time,
		}
	}

	return result, nil
}

func (s *Storage) CreatePayment(ctx context.Context, payment paymentmodel.Payment) (paymentmodel.Payment, error) {
	result, err := s.sqlc.CreatePayment(ctx, sqlc.CreatePaymentParams{
		AccountID: payment.AccountID,
		Method:    sqlc.PaymentMethod(payment.Method),
		Status:    sqlc.PaymentStatus(payment.Status),
		Total:     payment.Total.Int64(),
	})
	if err != nil {
		return paymentmodel.Payment{}, err
	}

	return paymentmodel.Payment{
		ID:          result.ID,
		AccountID:   result.AccountID,
		Method:      paymentmodel.PaymentMethod(result.Method),
		Status:      paymentmodel.PaymentStatus(result.Status),
		Total:       commonmodel.Concurrency(result.Total),
		DateCreated: result.DateCreated.Time,
	}, nil
}

type UpdatePaymentParams struct {
	ID     int64
	Method *paymentmodel.PaymentMethod
	Status *paymentmodel.PaymentStatus
	Total  *int64
}

func (s *Storage) UpdatePayment(ctx context.Context, params UpdatePaymentParams) (paymentmodel.Payment, error) {
	row, err := s.sqlc.UpdatePayment(ctx, sqlc.UpdatePaymentParams{
		ID:     params.ID,
		Method: *pgxptr.PtrBrandedToPgType(&sqlc.NullPaymentMethod{}, params.Method),
		Status: *pgxptr.PtrBrandedToPgType(&sqlc.NullPaymentStatus{}, params.Status),
		Total:  *pgxptr.PtrToPgtype(&pgtype.Int8{}, params.Total),
	})
	if err != nil {
		return paymentmodel.Payment{}, err
	}

	return paymentmodel.Payment{
		ID:          row.ID,
		AccountID:   row.AccountID,
		Method:      paymentmodel.PaymentMethod(row.Method),
		Status:      paymentmodel.PaymentStatus(row.Status),
		Total:       commonmodel.Concurrency(row.Total),
		DateCreated: row.DateCreated.Time,
	}, nil
}

func (s *Storage) DeletePayment(ctx context.Context, id int64) error {
	return s.sqlc.DeletePayment(ctx, id)
}

func (s *Storage) CreatePaymentItem(ctx context.Context, item paymentmodel.PaymentItem) (paymentmodel.PaymentItem, error) {
	row, err := s.sqlc.CreatePaymentItem(ctx, sqlc.CreatePaymentItemParams{
		PaymentID: item.PaymentID,
		Name:      item.Name,
		Price:     item.Price.Int64(),
	})
	if err != nil {
		return paymentmodel.PaymentItem{}, err
	}

	return paymentmodel.PaymentItem{
		ID:        row.ID,
		PaymentID: row.PaymentID,
		Name:      row.Name,
		Price:     commonmodel.Concurrency(row.Price),
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
