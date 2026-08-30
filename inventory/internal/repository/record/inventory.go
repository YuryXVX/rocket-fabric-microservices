package record

import (
	"time"

	"github.com/google/uuid"
)

type PartType int

const (
	// PartTypeUnspecified — не указан (дефолтное значение)
	PartTypeUnspecified PartType = iota // 0
	// PartTypeHull — корпус корабля
	PartTypeHull // 1
	// PartTypeEngine — двигатель
	PartTypeEngine // 2
	// PartTypeShield — защитный щит
	PartTypeShield // 3
	// PartTypeWeapon — вооружение
	PartTypeWeapon // 4
)

type Part struct {
	UUID          uuid.UUID
	Name          string
	Description   string
	Price         int64
	PartType      PartType
	StockQuantity int
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}
