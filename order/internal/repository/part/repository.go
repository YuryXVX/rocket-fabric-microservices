package part

import (
	"order/internal/repository/record"
	"sync"

	"github.com/google/uuid"
)

type repository struct {
	mu    sync.RWMutex
	store map[uuid.UUID][]record.PartRecord
}

func New() *repository {
	return &repository{
		store: make(map[uuid.UUID][]record.PartRecord),
	}
}
