package input

import (
	"github.com/google/uuid"
	"inventory/internal/model"
)

type PartFilter struct {
	// UUIDs — если не пустой, возвращаются только эти детали (приоритет)
	UUIDs []uuid.UUID
	// PartType — фильтр по типу (игнорируется если UUIDs заполнен)
	PartType model.PartType
}
