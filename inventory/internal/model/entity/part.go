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

func (p *part) UUID() uuid.UUID { return p.uuid }

func (p *part) Name() string { return p.name }

func (p *part) Reserved() int64 { return p.reserved }

func (p *part) Available() int64 { return p.stockQuantity - p.reserved }
