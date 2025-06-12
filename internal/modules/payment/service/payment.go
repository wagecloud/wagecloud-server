package paymentsvc

import (
	"context"
	"errors"

	paymentmodel "github.com/wagecloud/wagecloud-server/internal/modules/payment/model"
	paymentstorage "github.com/wagecloud/wagecloud-server/internal/modules/payment/storage"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
)

var (
	ErrPaymentNotFound = errors.New("payment not found")
	ErrInvalidPayment  = errors.New("invalid payment")
)

type Service interface {
	GetPayment(ctx context.Context, id int64) (paymentmodel.PaymentBase, error)
	ListPayments(ctx context.Context, params paymentstorage.ListPaymentsParams) (pagination.PaginateResult[paymentmodel.PaymentBase], error)
	CreatePayment(ctx context.Context, payment paymentmodel.PaymentBase) (paymentmodel.PaymentBase, error)
	UpdatePayment(ctx context.Context, params paymentstorage.UpdatePaymentParams) (paymentmodel.PaymentBase, error)
	DeletePayment(ctx context.Context, id int64) error
	CreatePaymentItem(ctx context.Context, item paymentmodel.PaymentItem) (paymentmodel.PaymentItem, error)
	CreatePaymentVNPAY(ctx context.Context, vnpay paymentmodel.PaymentVNPAY) (paymentmodel.PaymentVNPAY, error)
}

type ServiceImpl struct {
	storage *paymentstorage.Storage
}

func NewService(storage *paymentstorage.Storage) Service {
	return &ServiceImpl{storage: storage}
}

func (s *ServiceImpl) GetPayment(ctx context.Context, id int64) (paymentmodel.PaymentBase, error) {
	payment, err := s.storage.GetPayment(ctx, id)
	if err != nil {
		return paymentmodel.PaymentBase{}, err
	}
	return payment, nil
}

type ListPaymentsParams = paymentstorage.ListPaymentsParams

func (s *ServiceImpl) ListPayments(ctx context.Context, params ListPaymentsParams) (result pagination.PaginateResult[paymentmodel.PaymentBase], err error) {
	total, err := s.storage.CountPayments(ctx, params)
	if err != nil {
		return result, err
	}

	payments, err := s.storage.ListPayments(ctx, params)
	if err != nil {
		return result, err
	}

	return pagination.PaginateResult[paymentmodel.PaymentBase]{
		Total:    total,
		Data:     payments,
		Page:     params.Page,
		Limit:    params.Limit,
		NextPage: params.NextPage(total),
	}, nil
}

func (s *ServiceImpl) CreatePayment(ctx context.Context, payment paymentmodel.PaymentBase) (paymentmodel.PaymentBase, error) {
	return s.storage.CreatePayment(ctx, payment)
}

type UpdatePaymentParams = paymentstorage.UpdatePaymentParams

func (s *ServiceImpl) UpdatePayment(ctx context.Context, params UpdatePaymentParams) (paymentmodel.PaymentBase, error) {
	return s.storage.UpdatePayment(ctx, params)
}

func (s *ServiceImpl) DeletePayment(ctx context.Context, id int64) error {
	return s.storage.DeletePayment(ctx, id)
}

func (s *ServiceImpl) CreatePaymentItem(ctx context.Context, item paymentmodel.PaymentItem) (paymentmodel.PaymentItem, error) {
	return s.storage.CreatePaymentItem(ctx, item)
}

func (s *ServiceImpl) CreatePaymentVNPAY(ctx context.Context, vnpay paymentmodel.PaymentVNPAY) (paymentmodel.PaymentVNPAY, error) {
	return s.storage.CreatePaymentVNPAY(ctx, vnpay)
}
