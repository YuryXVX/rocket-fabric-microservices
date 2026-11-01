package api_iam

import "context"

type IAMService interface {
	Login(context.Context) error
	Logout(context.Context) error
	WhoIam(context.Context) error
}
