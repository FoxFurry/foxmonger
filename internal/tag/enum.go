package tag

import (
	"fmt"
	"math/rand"
	"strings"
)

type EnumProducer struct {
	enumValues []string
}

func (gen *EnumProducer) Produce() string {
	return gen.enumValues[rand.Intn(len(gen.enumValues))]
}

func (gen *EnumProducer) Initialize(input string) error {
	enumElements := strings.Split(input, modifierSplitter)

	if len(enumElements) != 2 { // Valid enum should contain 2 elements: enum keyword and enum values
		return fmt.Errorf("tag does not contain enum values: %s\n", input)
	}

	gen.enumValues = strings.Split(enumElements[1], enumValueSplitter)
	return nil
}

func NewEnumProducer(tag string) (Producer, error) {
	prod := &EnumProducer{}
	if err := prod.Initialize(tag); err != nil {
		return nil, err
	}

	return prod, nil
}
