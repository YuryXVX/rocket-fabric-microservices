package entity

import (
	errs "inventory/internal/errors"
	"inventory/internal/model/valueobject"
	"time"

	"github.com/google/uuid"
)

type Part = part

type part struct {
	uuid          uuid.UUID
	name          string
	description   string
	price         int64
	partType      *valueobject.PartType
	stockQuantity int64
	properties    *valueobject.PartProperties
	reserved      int64
	createdAt     time.Time
	updatedAt     *time.Time
}

func (p *Part) ReservedPart() error {
	if p.Available() <= 0 {
		return errs.ErrOutOfStock
	}

	p.reserved++

	return nil
}

func (p *Part) ReleasePart() error {
	if p.reserved <= 0 {
		return errs.ErrNothingToRelease
	}

	p.reserved--

	return nil
}

func RestorePart(partUUID uuid.UUID, name, description string, partType *valueobject.PartType, price int64,
	stockQuantity, reserved int64, properties *valueobject.PartProperties, createdAt time.Time) Part {
	return Part{
		uuid:          partUUID,
		name:          name,
		description:   description,
		partType:      partType,
		price:         price,
		stockQuantity: stockQuantity,
		reserved:      reserved,
		properties:    properties,
		createdAt:     createdAt,
	}
}

func (p *part) UUID() uuid.UUID { return p.uuid }

func (p *part) Name() string { return p.name }

func (p *part) Description() string { return p.description }

func (p *part) PartType() *valueobject.PartType { return p.partType }

func (p *part) Price() int64 { return p.price }

func (p *part) StockQuantity() int64 { return p.stockQuantity }

func (p *part) CreatedAt() time.Time { return p.createdAt }

func (p *part) Reserved() int64 { return p.reserved }

func (p *part) Available() int64 { return p.stockQuantity - p.reserved }
