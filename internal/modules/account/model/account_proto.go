package accountmodel

import accountv1 "github.com/wagecloud/wagecloud-server/gen/pb/account/v1"

func RoleModelToProto(role Role) accountv1.Role {
	return accountv1.Role(accountv1.Role_value[string(role)])
}

func RoleProtoToModel(role accountv1.Role) Role {
	return Role(accountv1.Role_name[int32(role)])
}

func AuthenticatedAccountProtoToModel(proto *accountv1.AuthenticatedAccount) AuthenticatedAccount {
	return AuthenticatedAccount{
		AccountID: proto.AccountId,
		Role:      RoleProtoToModel(proto.Role),
	}
}

func AuthenticatedAccountModelToProto(model AuthenticatedAccount) *accountv1.AuthenticatedAccount {
	return &accountv1.AuthenticatedAccount{
		AccountId: model.AccountID,
		Role:      RoleModelToProto(model.Role),
	}
}

func AccountUserProtoToModel(proto *accountv1.Account) AccountUser {
	return AccountUser{
		AccountBase: AccountBase{
			ID:        proto.Id,
			Role:      RoleProtoToModel(proto.Role),
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
		Role:      RoleModelToProto(model.Role),
		Username:  model.Username,
		Email:     model.Email,
		Name:      model.Name,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
