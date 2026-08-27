package record

import "github.com/google/uuid"

type PartRecord struct {
	PartUUID uuid.UUID
	PartType PartType
	Price    int64
}

type PartType string

const (
	PartTypeUnspecified PartType = "UNSPECIFIED"
	PartTypeHull        PartType = "HULL"
	PartTypeEngine      PartType = "ENGINE"
	PartTypeShield      PartType = "SHIELD"
	PartTypeWeapon      PartType = "WEAPON"
)
