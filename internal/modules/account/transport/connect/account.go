package accountconnect

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	accountv1 "github.com/wagecloud/wagecloud-server/gen/pb/account/v1"
	"github.com/wagecloud/wagecloud-server/gen/pb/account/v1/accountv1connect"
	accountmodel "github.com/wagecloud/wagecloud-server/internal/modules/account/model"
	accountsvc "github.com/wagecloud/wagecloud-server/internal/modules/account/service"
)

type ImplementedAccountServiceHandler struct {
	accountv1connect.UnimplementedAccountServiceHandler
	service accountsvc.Service
}

func NewAccountServiceHandler(service accountsvc.Service) (string, http.Handler) {
	return accountv1connect.NewAccountServiceHandler(&ImplementedAccountServiceHandler{
		service: service,
	})
}

func (t *ImplementedAccountServiceHandler) GetUser(ctx context.Context, req *connect.Request[accountv1.GetUserRequest]) (*connect.Response[accountv1.GetUserResponse], error) {
	result, err := t.service.GetUser(ctx, accountsvc.GetUserParams{})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&accountv1.GetUserResponse{
		Account: accountmodel.AccountUserModelToProto(result),
	}), err
}

func (t *ImplementedAccountServiceHandler) Login(ctx context.Context, req *connect.Request[accountv1.LoginRequest]) (*connect.Response[accountv1.LoginResponse], error) {
	result, err := t.service.LoginUser(context.Background(), accountsvc.LoginUserParams{
		Username: req.Msg.Username,
		Password: req.Msg.Password,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&accountv1.LoginResponse{
		Token:   result.Token,
		Account: accountmodel.AccountUserModelToProto(result.Account),
	}), nil
}

func (t *ImplementedAccountServiceHandler) Register(ctx context.Context, req *connect.Request[accountv1.RegisterRequest]) (*connect.Response[accountv1.RegisterResponse], error) {
	result, err := t.service.RegisterUser(context.Background(), accountsvc.RegisterUserParams{
		Name:     req.Msg.Name,
		Username: req.Msg.Username,
		Email:    req.Msg.Email,
		Password: req.Msg.Password,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&accountv1.RegisterResponse{
		Token:   result.Token,
		Account: accountmodel.AccountUserModelToProto(result.Account),
	}), nil
}
