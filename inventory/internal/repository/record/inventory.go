package record

import (
	"time"

	"github.com/google/uuid"
)

type Part struct {
	UUID          uuid.UUID
	Name          string
	Description   string
	PartType      string
	Price         int64
	StockQuantity int64
	Reserved      int64
	Properties    []byte
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}
