package accountmodel

import accountv1 "github.com/wagecloud/wagecloud-server/gen/pb/account/v1"

func AccountTypeModelToProto(accountType AccountType) accountv1.AccountType {
	return accountv1.AccountType(accountv1.AccountType_value[string(accountType)])
}

func AccountTypeProtoToModel(accountType accountv1.AccountType) AccountType {
	return AccountType(accountv1.AccountType_name[int32(accountType)])
}

func AuthenticatedAccountProtoToModel(proto *accountv1.AuthenticatedAccount) AuthenticatedAccount {
	return AuthenticatedAccount{
		AccountID: proto.AccountId,
		Type:      AccountTypeProtoToModel(proto.Type),
	}
}

func AuthenticatedAccountModelToProto(model AuthenticatedAccount) *accountv1.AuthenticatedAccount {
	return &accountv1.AuthenticatedAccount{
		AccountId: model.AccountID,
		Type:      AccountTypeModelToProto(model.Type),
	}
}

func AccountUserProtoToModel(proto *accountv1.Account) AccountUser {
	return AccountUser{
		AccountBase: AccountBase{
			ID:        proto.Id,
			Type:      AccountTypeProtoToModel(proto.Type),
			Username:  proto.Username,
			CreatedAt: proto.CreatedAt,
			UpdatedAt: proto.UpdatedAt,
			Name:      proto.Name,
		},
		Email: proto.Email,
	}
}

func AccountUserModelToProto(model AccountUser) *accountv1.Account {
	return &accountv1.Account{
		Id:        model.ID,
		Type:      AccountTypeModelToProto(model.Type),
		Username:  model.Username,
		Email:     model.Email,
		Name:      model.Name,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
