// Package tag exposes Producer and Modifier interfaces describing row generation. Implementations of these interfaces
// will correspond to their tags
package tag

import (
	"fmt"
)

const (
	ModifierSplitter  = ":"
	EnumValueSplitter = ","
)

type Tag interface {
	Initialize(string) error
}

type Producer interface {
	Tag
	Produce() string
}

type Modifier interface {
	Tag
	Modify(string) string
}

type Generator struct {
	producer  Producer
	modifiers []Modifier
	RowName   string
}

func (t *Generator) Do() string {
	source := t.producer.Produce()

	for idx := range t.modifiers {
		source = t.modifiers[idx].Modify(source)
	}

	return source
}

func (t *Generator) SetProducer(newProducer Producer) error {
	if t.producer != nil {
		return fmt.Errorf("producer already exists")
	}

	t.producer = newProducer
	return nil
}

func (t *Generator) AddModifier(newModifier Modifier) {
	t.modifiers = append(t.modifiers, newModifier)
}
