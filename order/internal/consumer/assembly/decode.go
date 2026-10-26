package assembly

import (
	"fmt"
	"order/internal/model"
	v1 "shared/pkg/proto/events/v1"

	"google.golang.org/protobuf/proto"
)

func decodeAssemblyRecorded(data []byte) (model.ShipAssembledEvent, error) {
	var pb v1.ShipAssembled
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.ShipAssembledEvent{}, fmt.Errorf("не удалось десериализовать protobuf: %w", err)
	}

	return model.ShipAssembledEvent{
		UserUuid:     pb.UserUuid,
		EventUuid:    pb.EventUuid,
		OrderUuid:    pb.OrderUuid,
		BuildTimeSec: pb.BuildTimeSec,
		AssembledAt:  pb.AssembledAt.GetSeconds(),
	}, nil
}
