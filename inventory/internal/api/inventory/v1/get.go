package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"inventory/internal/api/converter"
	errs "inventory/internal/errors"
	v1 "shared/pkg/proto/inventory/v1"
)

func (a *api) GetPartByUUID(ctx context.Context, req *v1.GetPartByUUIDRequest) (*v1.GetPartByUUIDResponse, error) {
	if req.GetUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid обязателен")
	}

	part, err := a.serviceInventory.Get(ctx, req.Uuid)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidUUID) {
			return nil, status.Errorf(codes.InvalidArgument, "UUID %s не правильного формата", req.GetUuid())
		}

		if errors.Is(err, errs.ErrPartNotFound) {
			return nil, status.Errorf(codes.NotFound, "деталь с UUID %s не найдена", req.GetUuid())
		}

		return nil, status.Errorf(codes.Internal, "ошибка получения детали")
	}

	return &v1.GetPartByUUIDResponse{
		Part: converter.ModelPartToProtoPart(part),
	}, nil
}
