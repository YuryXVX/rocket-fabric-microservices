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
	class EngineClass
}

func NewEngineProperties(engineClass EngineClass) (*PartProperties, error) {
	if engineClass == "" {
		return nil, fmt.Errorf("engine class cannot be empty")
	}

	switch engineClass {
	case EngineClassA, EngineClassB, EngineClassC:
	default:
		return nil, fmt.Errorf("invalid engine class: %q", engineClass)
	}

	return &PartProperties{
		engine: &engineProperties{
			class: engineClass,
		},
	}, nil
}

func (e *engineProperties) Class() EngineClass { return e.class }
