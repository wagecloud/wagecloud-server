package accountmodel

import commonmodel "github.com/wagecloud/wagecloud-server/internal/shared/model"

var (
	ErrAccountNotFound      = commonmodel.NewError("ErrAccountNotFound", "Account not found")
	ErrAccountAlreadyExists = commonmodel.NewError("ErrAccountAlreadyExists", "Account already exists")
	ErrInvalidCredentials   = commonmodel.NewError("ErrInvalidCredentials", "Invalid credentials provided")
	ErrWrongCurrentPassword = commonmodel.NewError("ErrWrongCurrentPassword", "Wrong current password provided")
)
