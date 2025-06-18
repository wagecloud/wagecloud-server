package paymentconnect

import (
	"context"

	"connectrpc.com/connect"
	paymentv1 "github.com/wagecloud/wagecloud-server/gen/pb/payment/v1"
	"github.com/wagecloud/wagecloud-server/gen/pb/payment/v1/paymentv1connect"
	paymentmodel "github.com/wagecloud/wagecloud-server/internal/modules/payment/model"
	paymentservice "github.com/wagecloud/wagecloud-server/internal/modules/payment/service"
	paymentstorage "github.com/wagecloud/wagecloud-server/internal/modules/payment/storage"
	commonmodel "github.com/wagecloud/wagecloud-server/internal/shared/model"
	"github.com/wagecloud/wagecloud-server/internal/utils/ptr"
	"github.com/wagecloud/wagecloud-server/internal/utils/slice"
)

type ImplementedPaymentServiceHandler struct {
	paymentv1connect.UnimplementedPaymentServiceHandler
	service paymentservice.Service
}

func NewImplementedPaymentServiceHandler(service paymentservice.Service) paymentv1connect.PaymentServiceHandler {
	return &ImplementedPaymentServiceHandler{
		service: service,
	}
}

func (t *ImplementedPaymentServiceHandler) GetPayment(ctx context.Context, req *connect.Request[paymentv1.GetPaymentRequest]) (*connect.Response[paymentv1.GetPaymentResponse], error) {
	result, err := t.service.GetPayment(ctx, req.Msg.Id)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&paymentv1.GetPaymentResponse{
		Payment: paymentmodel.PaymentModelToProto(result),
	}), nil
}

func (t *ImplementedPaymentServiceHandler) ListPayments(ctx context.Context, req *connect.Request[paymentv1.ListPaymentsRequest]) (*connect.Response[paymentv1.ListPaymentsResponse], error) {
	result, err := t.service.ListPayments(ctx, paymentstorage.ListPaymentsParams{
		PaginationParams: commonmodel.PaginationParamsProtoToModel(req.Msg.Pagination),
		AccountID:        req.Msg.AccountId,
		Method:           ptr.Convert(req.Msg.Method, paymentmodel.PaymentMethodProtoToModel),
		Status:           ptr.Convert(req.Msg.Status, paymentmodel.PaymentStatusProtoToModel),
		DateCreatedFrom:  req.Msg.DateCreatedFrom,
		DateCreatedTo:    req.Msg.DateCreatedTo,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&paymentv1.ListPaymentsResponse{
		Payments:   slice.Map(result.Data, paymentmodel.PaymentModelToProto),
		Pagination: commonmodel.PaginateResultModelToProto(result),
	}), nil
}

// func (t *ImplementedPaymentServiceHandler) CreatePayment(ctx context.Context, req *connect.Request[paymentv1.CreatePaymentRequest]) (*connect.Response[paymentv1.CreatePaymentResponse], error) {
// 	result, err := t.service.CreatePayment(ctx, paymentmodel.Payment{
// 		Method: paymentmodel.PaymentMethodProtoToModel(req.Msg.Method),
// 		Total:  int64(req.Msg.Total),
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return connect.NewResponse(&paymentv1.CreatePaymentResponse{
// 		Payment: paymentmodel.PaymentModelToProto(result),
// 	}), nil
// }

func (t *ImplementedPaymentServiceHandler) UpdatePayment(ctx context.Context, req *connect.Request[paymentv1.UpdatePaymentRequest]) (*connect.Response[paymentv1.UpdatePaymentResponse], error) {
	_, err := t.service.UpdatePayment(ctx, paymentstorage.UpdatePaymentParams{
		ID:     req.Msg.Id,
		Method: ptr.Convert(req.Msg.Method, paymentmodel.PaymentMethodProtoToModel),
		Status: ptr.Convert(req.Msg.Status, paymentmodel.PaymentStatusProtoToModel),
		Total:  req.Msg.Total,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&paymentv1.UpdatePaymentResponse{}), nil
}

func (t *ImplementedPaymentServiceHandler) DeletePayment(ctx context.Context, req *connect.Request[paymentv1.DeletePaymentRequest]) (*connect.Response[paymentv1.DeletePaymentResponse], error) {
	err := t.service.DeletePayment(ctx, req.Msg.Id)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&paymentv1.DeletePaymentResponse{}), nil
}

// func (t *ImplementedPaymentServiceHandler) CreatePaymentItem(ctx context.Context, req *connect.Request[paymentv1.CreatePaymentItemRequest]) (*connect.Response[paymentv1.CreatePaymentItemResponse], error) {
// 	result, err := t.service.CreatePaymentItem(ctx, paymentmodel.PaymentItem{
// 		PaymentID: req.Msg.PaymentId,
// 		Name:      req.Msg.Name,
// 		Price:     req.Msg.Price,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return connect.NewResponse(&paymentv1.CreatePaymentItemResponse{
// 		PaymentItem: paymentmodel.PaymentItemModelToProto(result),
// 	}), nil
// }

// func (t *ImplementedPaymentServiceHandler) CreateVNPAYPayment(ctx context.Context, req *connect.Request[paymentv1.CreateVNPAYPaymentRequest]) (*connect.Response[paymentv1.CreateVNPAYPaymentResponse], error) {
// 	result, err := t.service.CreatePaymentVNPAY(ctx, paymentmodel.PaymentVNPAY{
// 		ID:                 req.Msg.Id,
// 		VnpTxnRef:          req.Msg.VnpTxnRef,
// 		VnpOrderInfo:       req.Msg.VnpOrderInfo,
// 		VnpTransactionNo:   req.Msg.VnpTransactionNo,
// 		VnpTransactionDate: req.Msg.VnpTransactionDate,
// 		VnpCreateDate:      req.Msg.VnpCreateDate,
// 		VnpIpAddr:          req.Msg.VnpIpAddr,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return connect.NewResponse(&paymentv1.CreateVNPAYPaymentResponse{
// 		VnpayPayment: paymentmodel.VnpayPaymentModelToProto(result),
// 	}), nil
// }
