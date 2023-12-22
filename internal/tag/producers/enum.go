package producers

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/FoxFurry/foxmonger/internal/tag"
)

type EnumProducer struct {
	enumValues []string
}

func (gen *EnumProducer) Produce() string {
	return gen.enumValues[rand.Intn(len(gen.enumValues))]
}

func (gen *EnumProducer) Initialize(input string) error {
	enumElements := strings.Split(input, tag.ModifierSplitter)

	if len(enumElements) != 2 { // Valid enum should contain 2 elements: enum keyword and enum values
		return fmt.Errorf("tag does not contain enum values: %s\n", input)
	}

	gen.enumValues = strings.Split(enumElements[1], tag.EnumValueSplitter)
	return nil
}

func NewEnumProducer(tag string) (tag.Producer, error) {
	prod := &EnumProducer{}
	if err := prod.Initialize(tag); err != nil {
		return nil, err
	}

	return prod, nil
}
