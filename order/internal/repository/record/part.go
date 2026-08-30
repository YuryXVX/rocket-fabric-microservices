package record

import "github.com/google/uuid"

type PartRecord struct {
	OrderUUID uuid.UUID `db:"order_uuid"`
	PartUUID  uuid.UUID `db:"part_uuid"`
	PartType  PartType  `db:"part_type"`
	Price     int64     `db:"price"`
}

type PartType string

const (
	PartTypeUnspecified PartType = "UNSPECIFIED"
	PartTypeHull        PartType = "HULL"
	PartTypeEngine      PartType = "ENGINE"
	PartTypeShield      PartType = "SHIELD"
	PartTypeWeapon      PartType = "WEAPON"
)
