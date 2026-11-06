package interceptor_auth

import (
	"context"
	"inventory/internal/model"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type IamClient interface {
	WhoIam(ctx context.Context, sessionUUID string) (model.Session, error)
}

const authorizationHeader = "Authorization"

type ctxKey int

const tokenKey ctxKey = 0

func Auth(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
	iamClient IamClient,
) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "отсутствуют metadata")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, status.Error(codes.Unauthenticated, "отсутствует authorization header")
	}

	if len(values) == 0 {
		return nil, status.Error(codes.Unauthenticated, "<token> отсутвует")
	}

	token := values[0]

	if token == "" {
		return nil, status.Error(codes.Unauthenticated, "формат: <token>")
	}

	ctx = context.WithValue(ctx, tokenKey, token)

	_, err := iamClient.WhoIam(ctx, token)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "сессия пользоватля не активна")
	}

	return handler(ctx, req)
}

// TokenFromContext возвращает токен из контекста
func TokenFromContext(ctx context.Context) string {
	token, _ := ctx.Value(tokenKey).(string)
	return token
}

func AuthorizationInterceptor(iamClient IamClient) func(ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		return Auth(ctx, req, info, handler, iamClient)
	}
}
