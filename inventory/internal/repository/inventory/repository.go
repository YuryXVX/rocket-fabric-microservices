package inventory

import (
	"sync"

	"inventory/internal/repository/record"
)

type repository struct {
	mu    sync.RWMutex
	parts map[string]record.Part
}

func NewRepository() *repository {
	return &repository{
		parts: createMocks(),
	}
}
