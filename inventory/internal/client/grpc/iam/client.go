package iam

import (
	"context"
	"inventory/internal/model"
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

func (c *IamClient) WhoIam(ctx context.Context, sessionUUID string) (model.Session, error) {
	_, err := c.grpc.Whoami(ctx, &iamV1.WhoamiRequest{
		SessionUuid: sessionUUID,
	})

	if err != nil {
		return model.Session{}, err
	}

	return model.Session{}, nil
}
