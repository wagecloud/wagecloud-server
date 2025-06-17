package accountsvc

// type ServiceRpcImpl struct {
// 	connect accountv1connect.AccountServiceClient
// }

// func NewServiceRpc(connect accountv1connect.AccountServiceClient) Service {
// 	return &ServiceRpcImpl{
// 		connect: connect,
// 	}
// }

// func (s *ServiceRpcImpl) GetUser(ctx context.Context, params GetUserParams) (accountmodel.AccountUser, error) {
// 	result, err := s.connect.GetUser(ctx, connect.NewRequest(&accountv1.GetUserRequest{
// 		Id:       params.ID,
// 		Username: params.Username,
// 		Email:    params.Email,
// 	}))
// 	if err != nil {
// 		return accountmodel.AccountUser{}, err
// 	}

// 	return accountmodel.AccountUserProtoToModel(result.Msg.Account), nil
// }

// func (s *ServiceRpcImpl) LoginUser(ctx context.Context, params LoginUserParams) (LoginUserResult, error) {
// 	result, err := s.connect.Login(ctx, connect.NewRequest(&accountv1.LoginRequest{
// 		Username: params.Username,
// 		Password: params.Password,
// 	}))
// 	if err != nil {
// 		return LoginUserResult{}, err
// 	}

// 	return LoginUserResult{
// 		Token:   result.Msg.Token,
// 		Account: accountmodel.AccountUserProtoToModel(result.Msg.Account),
// 	}, nil
// }

// func (s *ServiceRpcImpl) RegisterUser(ctx context.Context, params RegisterUserParams) (RegisterUserResult, error) {
// 	result, err := s.connect.Register(ctx, connect.NewRequest(&accountv1.RegisterRequest{
// 		Name:     params.Name,
// 		Username: params.Username,
// 		Email:    params.Email,
// 		Password: params.Password,
// 	}))
// 	if err != nil {
// 		return RegisterUserResult{}, err
// 	}

// 	return RegisterUserResult{
// 		Token:   result.Msg.Token,
// 		Account: accountmodel.AccountUserProtoToModel(result.Msg.Account),
// 	}, nil
// }
