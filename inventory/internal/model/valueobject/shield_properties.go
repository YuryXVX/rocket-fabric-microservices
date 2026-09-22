package valueobject

import "fmt"

type ShieldProperties = shieldProperties

type ShieldType = string

const (
	ShieldPlasma ShieldType = "plasma"
	ShieldEnergy ShieldType = "energy"
)

type shieldProperties struct {
	shieldType ShieldType
}

func NewShieldProperties(shieldType ShieldType) (*PartProperties, error) {
	switch shieldType {
	case ShieldPlasma, ShieldEnergy:
	default:
		return nil, fmt.Errorf("invalid shield type: %q", shieldType)
	}

	return &PartProperties{
		shield: &ShieldProperties{
			shieldType: shieldType,
		},
	}, nil
}
