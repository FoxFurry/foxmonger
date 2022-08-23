package main

import (
	"fmt"
)

type Generator interface {
	Generate() string
}

type Modifier interface {
	Modify(string) string
}

type tagGenerator struct {
	generator Generator
	modifiers []Modifier
}

func (r *tagGenerator) do() string {
	source := r.generator.Generate()

	for idx := range r.modifiers {
		source = r.modifiers[idx].Modify(source)
	}

	return source
}

func (r *tagGenerator) setGenerator(generator Generator) {
	if r.generator != nil {
		fmt.Println("generator is already set, ignoring succeeding generators")
	} else {
		r.generator = generator
	}
}

func (r *tagGenerator) addModifier(modifier Modifier) {
	r.modifiers = append(r.modifiers, modifier)
}
