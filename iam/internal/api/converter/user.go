package converter

import (
	"iam/internal/service/input"
	v1 "shared/pkg/proto/user/v1"
)

func RequestToInput(req *v1.RegisterRequest) input.RegisterInput {
	return input.RegisterInput{
		Login:    req.Info.Info.GetLogin(),
		Password: req.Info.GetPassword(),
	}
}
