package domain

import (
	errs "inventory/internal/errors"
	"inventory/internal/model/entity"
	"inventory/internal/model/valueobject"
)

type partsSet struct {
	hull   *valueobject.HullProperties
	engine *valueobject.EngineProperties
	shield *valueobject.ShieldProperties
	weapon *valueobject.WeaponProperties
}

type CompatibilityChecker struct{}

func New() *CompatibilityChecker {
	return &CompatibilityChecker{}
}

func (c *CompatibilityChecker) Check(parts []entity.Part) error {
	ps := c.extractParts(parts)

	if err := c.requirementDetailCheck(ps); err != nil {
		return err
	}

	if err := c.checkHull(ps); err != nil {
		return err
	}

	if err := c.checkShield(ps); err != nil {
		return err
	}

	return nil
}

func (c *CompatibilityChecker) requirementDetailCheck(ps partsSet) error {
	if ps.hull == nil || ps.engine == nil {
		return errs.ErrRequirementDetailCheck
	}

	return nil
}

func (c *CompatibilityChecker) checkHull(ps partsSet) error {
	if !ps.hull.CanSupport(ps.engine) {
		return errs.ErrIncompatibleHull
	}

	return nil
}

func (c *CompatibilityChecker) checkShield(ps partsSet) error {
	if ps.shield == nil || ps.weapon == nil {
		return nil
	}

	switch {
	case ps.shield.ShieldType() == valueobject.ShieldEnergy && ps.weapon.WeaponType() == valueobject.WeaponLaserType,
		ps.shield.ShieldType() == valueobject.ShieldPlasma && ps.weapon.WeaponType() == valueobject.WeaponMissile,
		ps.shield.ShieldType() == valueobject.ShieldEnergy && ps.weapon.WeaponType() == valueobject.WeaponMissile:
		return nil
	default:
		return errs.ErrIncompatibleShield
	}
}

func (c *CompatibilityChecker) extractParts(parts []entity.Part) partsSet {
	var ps partsSet

	for _, part := range parts {
		props := part.Properties()

		if props.Engine() != nil {
			ps.engine = props.Engine()
		}

		if props.Hull() != nil {
			ps.hull = props.Hull()
		}

		if props.Shield() != nil {
			ps.shield = props.Shield()
		}

		if props.Weapon() != nil {
			ps.weapon = props.Weapon()
		}
	}

	return ps
}
