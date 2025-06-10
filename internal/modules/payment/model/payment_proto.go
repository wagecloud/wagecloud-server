package paymentmodel

import paymentv1 "github.com/wagecloud/wagecloud-server/gen/pb/payment/v1"

func PaymentModelToProto(payment PaymentBase) *paymentv1.Payment {
	return &paymentv1.Payment{
		Id:          payment.ID,
		AccountId:   payment.AccountID,
		Method:      PaymentMethodModelToProto(payment.Method),
		Status:      PaymentStatusModelToProto(payment.Status),
		Total:       payment.Total,
		DateCreated: payment.DateCreated,
	}
}

func PaymentProtoToModel(payment *paymentv1.Payment) PaymentBase {
	return PaymentBase{
		ID:          payment.Id,
		AccountID:   payment.AccountId,
		Method:      PaymentMethodProtoToModel(payment.Method),
		Status:      PaymentStatusProtoToModel(payment.Status),
		Total:       payment.Total,
		DateCreated: payment.DateCreated,
	}
}

func PaymentItemModelToProto(item PaymentItem) *paymentv1.PaymentItem {
	return &paymentv1.PaymentItem{
		Id:        item.ID,
		PaymentId: item.PaymentID,
		Name:      item.Name,
		Price:     item.Price,
	}
}

func PaymentItemProtoToModel(item *paymentv1.PaymentItem) PaymentItem {
	return PaymentItem{
		ID:        item.Id,
		PaymentID: item.PaymentId,
		Name:      item.Name,
		Price:     item.Price,
	}
}

func VnpayPaymentModelToProto(vnpay PaymentVNPAY) *paymentv1.VNPAYPayment {
	return &paymentv1.VNPAYPayment{
		Id:                 vnpay.ID,
		VnpTxnRef:          vnpay.VnpTxnRef,
		VnpOrderInfo:       vnpay.VnpOrderInfo,
		VnpTransactionNo:   vnpay.VnpTransactionNo,
		VnpTransactionDate: vnpay.VnpTransactionDate,
		VnpCreateDate:      vnpay.VnpCreateDate,
		VnpIpAddr:          vnpay.VnpIpAddr,
	}
}

func VnpayPaymentProtoToModel(vnpay *paymentv1.VNPAYPayment) PaymentVNPAY {
	return PaymentVNPAY{
		ID:                 vnpay.Id,
		VnpTxnRef:          vnpay.VnpTxnRef,
		VnpOrderInfo:       vnpay.VnpOrderInfo,
		VnpTransactionNo:   vnpay.VnpTransactionNo,
		VnpTransactionDate: vnpay.VnpTransactionDate,
		VnpCreateDate:      vnpay.VnpCreateDate,
		VnpIpAddr:          vnpay.VnpIpAddr,
	}
}

func PaymentMethodProtoToModel(method paymentv1.PaymentMethod) PaymentMethod {
	return PaymentMethod(method.String())
}

func PaymentMethodModelToProto(method PaymentMethod) paymentv1.PaymentMethod {
	return paymentv1.PaymentMethod(paymentv1.PaymentMethod_value[string(method)])
}

func PaymentStatusProtoToModel(status paymentv1.PaymentStatus) PaymentStatus {
	return PaymentStatus(status.String())
}

func PaymentStatusModelToProto(status PaymentStatus) paymentv1.PaymentStatus {
	return paymentv1.PaymentStatus(paymentv1.PaymentStatus_value[string(status)])
}
