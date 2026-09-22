package valueobject

import "fmt"

type EngineProperties = engineProperties

type EngineClass = string

const (
	EngineClassA EngineClass = "A"
	EngineClassB EngineClass = "B"
	EngineClassC EngineClass = "C"
)

type engineProperties struct {
	class            EngineClass
	requiredStrength int
}

func NewEngineProperties(engineClass EngineClass, requiredStrength int) (*PartProperties, error) {
	if engineClass == "" {
		return nil, fmt.Errorf("engine class cannot be empty")
	}

	if requiredStrength <= 0 {
		return nil, fmt.Errorf("required strength must be positive, got %d", requiredStrength)
	}

	switch engineClass {
	case EngineClassA, EngineClassB, EngineClassC:
	default:
		return nil, fmt.Errorf("invalid engine class: %q", engineClass)
	}

	return &PartProperties{
		engine: &engineProperties{
			class:            engineClass,
			requiredStrength: requiredStrength,
		},
	}, nil
}

func (e *engineProperties) Class() EngineClass    { return e.class }
func (e *engineProperties) RequiredStrength() int { return e.requiredStrength }
