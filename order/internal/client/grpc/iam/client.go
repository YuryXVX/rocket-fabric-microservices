package iam

import (
	"context"
	iamV1 "shared/pkg/proto/auth/v1"
)

type IamClient struct {
	grpc iamV1.AuthServiceClient
}

func New(grpc iamV1.AuthServiceClient) *IamClient {
	return &IamClient{
		grpc: grpc,
	}
}

func (c *IamClient) WhoIam(ctx context.Context, sessionUUID string) (string, error) {
	res, err := c.grpc.Whoami(ctx, &iamV1.WhoamiRequest{
		SessionUuid: sessionUUID,
	})

	if err != nil {
		return "", err
	}

	return res.User.Uuid, nil
}
