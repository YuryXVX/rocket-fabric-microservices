package valueobject

import "fmt"

type WeaponProperties = weaponProperties

type WeaponType = string

const (
	WeaponLaserType = "laser"
	WeaponMissile   = "missile"
)

type weaponProperties struct {
	weaponType WeaponType
}

func NewWeaponProperties(weaponType WeaponType) (*PartProperties, error) {
	switch weaponType {
	case WeaponLaserType, WeaponMissile:
	default:
		return nil, fmt.Errorf("invalid weapon type: %q", weaponType)
	}

	return &PartProperties{
		weapon: &weaponProperties{
			weaponType: weaponType,
		},
	}, nil
}
