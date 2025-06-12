package paymentsvc

import (
	"context"
	"errors"
	"strconv"

	accountmodel "github.com/wagecloud/wagecloud-server/internal/modules/account/model"
	paymentmodel "github.com/wagecloud/wagecloud-server/internal/modules/payment/model"
	paymentstorage "github.com/wagecloud/wagecloud-server/internal/modules/payment/storage"
	commonmodel "github.com/wagecloud/wagecloud-server/internal/shared/model"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
)

var (
	ErrPaymentNotFound = errors.New("payment not found")
	ErrInvalidPayment  = errors.New("invalid payment")
)

type Service interface {
	GetPayment(ctx context.Context, id int64) (paymentmodel.Payment, error)
	ListPayments(ctx context.Context, params ListPaymentsParams) (pagination.PaginateResult[paymentmodel.Payment], error)
	CreatePayment(ctx context.Context, params CreatePaymentParams) (CreatePaymentResult, error)
	UpdatePayment(ctx context.Context, params UpdatePaymentParams) (paymentmodel.Payment, error)
	DeletePayment(ctx context.Context, id int64) error
}

type ServiceImpl struct {
	storage   *paymentstorage.Storage
	platforms map[paymentmodel.PaymentMethod]PaymentPlatform
}

func NewService(storage *paymentstorage.Storage) Service {
	return &ServiceImpl{
		storage: storage,
		platforms: map[paymentmodel.PaymentMethod]PaymentPlatform{
			paymentmodel.PaymentMethodVNPAY: &VnpayPlatform{},
			// paymentmodel.PaymentMethodMOMO:  &MomoPlatform{},
		},
	}
}

func (s *ServiceImpl) GetPayment(ctx context.Context, id int64) (paymentmodel.Payment, error) {
	payment, err := s.storage.GetPayment(ctx, id)
	if err != nil {
		return paymentmodel.Payment{}, err
	}
	return payment, nil
}

type ListPaymentsParams = paymentstorage.ListPaymentsParams

func (s *ServiceImpl) ListPayments(ctx context.Context, params ListPaymentsParams) (result pagination.PaginateResult[paymentmodel.Payment], err error) {
	total, err := s.storage.CountPayments(ctx, params)
	if err != nil {
		return result, err
	}

	payments, err := s.storage.ListPayments(ctx, params)
	if err != nil {
		return result, err
	}

	return pagination.PaginateResult[paymentmodel.Payment]{
		Total:    total,
		Data:     payments,
		Page:     params.Page,
		Limit:    params.Limit,
		NextPage: params.NextPage(total),
	}, nil
}

type CreatePaymentParams struct {
	Account accountmodel.AuthenticatedAccount
	Method  paymentmodel.PaymentMethod
	Items   []CreatePaymentParamsItem
}

type CreatePaymentParamsItem struct {
	Name  string
	Price commonmodel.Concurrency
}

type CreatePaymentResult struct {
	Payment paymentmodel.Payment
	Items   []paymentmodel.PaymentItem
	URL     string
}

func (s *ServiceImpl) CreatePayment(ctx context.Context, params CreatePaymentParams) (CreatePaymentResult, error) {
	txStorage, err := s.storage.BeginTx(ctx)
	if err != nil {
		return CreatePaymentResult{}, err
	}
	defer txStorage.Rollback(ctx)

	var totalPrice commonmodel.Concurrency

	for _, item := range params.Items {
		totalPrice += item.Price
	}

	payment, err := txStorage.CreatePayment(ctx, paymentmodel.Payment{
		AccountID: params.Account.AccountID,
		Method:    params.Method,
		Status:    paymentmodel.PaymentStatusPending,
		Total:     totalPrice,
	})
	if err != nil {
		return CreatePaymentResult{}, err
	}

	var items []paymentmodel.PaymentItem
	for _, item := range params.Items {
		paymentItem, err := txStorage.CreatePaymentItem(ctx, paymentmodel.PaymentItem{
			PaymentID: payment.ID,
			Name:      item.Name,
			Price:     commonmodel.Concurrency(item.Price),
		})
		if err != nil {
			return CreatePaymentResult{}, err
		}

		items = append(items, paymentItem)
	}

	// Create payment url
	platform, ok := s.platforms[params.Method]
	if !ok {
		return CreatePaymentResult{}, ErrInvalidPayment
	}

	url, err := platform.CreateOrder(ctx, CreateOrderParams{
		PaymentID: payment.ID,
		Info:      "Payment for account " + strconv.FormatInt(params.Account.AccountID, 10),
		Amount:    totalPrice,
	})
	if err != nil {
		return CreatePaymentResult{}, err
	}

	return CreatePaymentResult{
		Payment: payment,
		Items:   items,
		URL:     url,
	}, txStorage.Commit(ctx)
}

type UpdatePaymentParams = paymentstorage.UpdatePaymentParams

func (s *ServiceImpl) UpdatePayment(ctx context.Context, params UpdatePaymentParams) (paymentmodel.Payment, error) {
	return s.storage.UpdatePayment(ctx, params)
}

func (s *ServiceImpl) DeletePayment(ctx context.Context, id int64) error {
	return s.storage.DeletePayment(ctx, id)
}
