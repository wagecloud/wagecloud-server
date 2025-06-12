package paymentsvc

import (
	"context"

	"connectrpc.com/connect"
	paymentv1 "github.com/wagecloud/wagecloud-server/gen/pb/payment/v1"
	"github.com/wagecloud/wagecloud-server/gen/pb/payment/v1/paymentv1connect"
	paymentmodel "github.com/wagecloud/wagecloud-server/internal/modules/payment/model"
	paymentstorage "github.com/wagecloud/wagecloud-server/internal/modules/payment/storage"
	commonmodel "github.com/wagecloud/wagecloud-server/internal/shared/model"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
	"github.com/wagecloud/wagecloud-server/internal/utils/ptr"
	"github.com/wagecloud/wagecloud-server/internal/utils/slice"
)

type ServiceRpcImpl struct {
	connect paymentv1connect.PaymentServiceClient
}

func NewServiceRpc(connect paymentv1connect.PaymentServiceClient) Service {
	return &ServiceRpcImpl{
		connect: connect,
	}
}

func (s *ServiceRpcImpl) GetPayment(ctx context.Context, id int64) (paymentmodel.PaymentBase, error) {
	result, err := s.connect.GetPayment(ctx, connect.NewRequest(&paymentv1.GetPaymentRequest{
		Id: id,
	}))
	if err != nil {
		return paymentmodel.PaymentBase{}, err
	}

	return paymentmodel.PaymentProtoToModel(result.Msg.Payment), nil
}

func (s *ServiceRpcImpl) ListPayments(ctx context.Context, params paymentstorage.ListPaymentsParams) (pagination.PaginateResult[paymentmodel.PaymentBase], error) {
	result, err := s.connect.ListPayments(ctx, connect.NewRequest(&paymentv1.ListPaymentsRequest{
		Pagination:      commonmodel.PaginationParamsModelToProto(params.PaginationParams),
		Id:              params.ID,
		AccountId:       params.AccountID,
		Method:          ptr.Convert(params.Method, paymentmodel.PaymentMethodModelToProto),
		Status:          ptr.Convert(params.Status, paymentmodel.PaymentStatusModelToProto),
		DateCreatedFrom: params.DateCreatedFrom,
		DateCreatedTo:   params.DateCreatedTo,
	}))
	if err != nil {
		return pagination.PaginateResult[paymentmodel.PaymentBase]{}, err
	}

	return pagination.PaginateResult[paymentmodel.PaymentBase]{
		Data:       slice.Map(result.Msg.Payments, paymentmodel.PaymentProtoToModel),
		Limit:      result.Msg.Pagination.Limit,
		Page:       result.Msg.Pagination.Page,
		Total:      result.Msg.Pagination.Total,
		NextPage:   result.Msg.Pagination.NextPage,
		NextCursor: result.Msg.Pagination.NextCursor,
	}, nil
}

func (s *ServiceRpcImpl) CreatePayment(ctx context.Context, payment paymentmodel.PaymentBase) (paymentmodel.PaymentBase, error) {
	result, err := s.connect.CreatePayment(ctx, connect.NewRequest(&paymentv1.CreatePaymentRequest{
		Method: paymentmodel.PaymentMethodModelToProto(payment.Method),
		Total:  payment.Total,
	}))
	if err != nil {
		return paymentmodel.PaymentBase{}, err
	}

	return paymentmodel.PaymentProtoToModel(result.Msg.Payment), nil
}

func (s *ServiceRpcImpl) UpdatePayment(ctx context.Context, params paymentstorage.UpdatePaymentParams) (paymentmodel.PaymentBase, error) {
	result, err := s.connect.UpdatePayment(ctx, connect.NewRequest(&paymentv1.UpdatePaymentRequest{
		Id:     params.ID,
		Method: ptr.Convert(params.Method, paymentmodel.PaymentMethodModelToProto),
		Status: ptr.Convert(params.Status, paymentmodel.PaymentStatusModelToProto),
		Total:  params.Total,
	}))
	if err != nil {
		return paymentmodel.PaymentBase{}, err
	}

	return paymentmodel.PaymentProtoToModel(result.Msg.Payment), nil
}

func (s *ServiceRpcImpl) DeletePayment(ctx context.Context, id int64) error {
	_, err := s.connect.DeletePayment(ctx, connect.NewRequest(&paymentv1.DeletePaymentRequest{
		Id: id,
	}))
	return err
}

func (s *ServiceRpcImpl) CreatePaymentItem(ctx context.Context, item paymentmodel.PaymentItem) (paymentmodel.PaymentItem, error) {
	result, err := s.connect.CreatePaymentItem(ctx, connect.NewRequest(&paymentv1.CreatePaymentItemRequest{
		PaymentId: item.PaymentID,
		Name:      item.Name,
		Price:     item.Price,
	}))
	if err != nil {
		return paymentmodel.PaymentItem{}, err
	}

	return paymentmodel.PaymentItemProtoToModel(result.Msg.PaymentItem), nil
}

func (s *ServiceRpcImpl) CreatePaymentVNPAY(ctx context.Context, vnpay paymentmodel.PaymentVNPAY) (paymentmodel.PaymentVNPAY, error) {
	result, err := s.connect.CreateVNPAYPayment(ctx, connect.NewRequest(&paymentv1.CreateVNPAYPaymentRequest{
		Id:                 vnpay.ID,
		VnpTxnRef:          vnpay.VnpTxnRef,
		VnpOrderInfo:       vnpay.VnpOrderInfo,
		VnpTransactionNo:   vnpay.VnpTransactionNo,
		VnpTransactionDate: vnpay.VnpTransactionDate,
		VnpCreateDate:      vnpay.VnpCreateDate,
		VnpIpAddr:          vnpay.VnpIpAddr,
	}))
	if err != nil {
		return paymentmodel.PaymentVNPAY{}, err
	}

	return paymentmodel.VnpayPaymentProtoToModel(result.Msg.VnpayPayment), nil
}
