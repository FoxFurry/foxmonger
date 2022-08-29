package main

import (
	"fmt"
)

type Producer interface {
	Produce() string
}

type Modifier interface {
	Modify(string) string
}

type Generator struct {
	producer  Producer
	modifiers []Modifier
}

func NewEmptyGenerator() Generator {
	return Generator{}
}

func (r *Generator) do() string {
	source := r.producer.Produce()

	for idx := range r.modifiers {
		source = r.modifiers[idx].Modify(source)
	}

	return source
}

func (r *Generator) setProducer(producer Producer) {
	if r.producer != nil {
		fmt.Println("generator is already set, ignoring succeeding generators")
	} else {
		r.producer = producer
	}
}

func (r *Generator) addModifier(modifier Modifier) {
	r.modifiers = append(r.modifiers, modifier)
}
