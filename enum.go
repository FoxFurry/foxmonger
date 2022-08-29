package main

import (
	"fmt"
	"math/rand"
	"strings"
)

type EnumProducer struct {
	enumTag string
}

func (gen *EnumProducer) Produce() string {
	enumElements := strings.Split(gen.enumTag, modifierSplitter)

	if len(enumElements) != 2 { // Valid enum should contain 2 elements: enum keyword and enum values
		fmt.Printf("ignoring invalid enum tag: %s\n", gen.enumTag)
		return ""
	}

	enumValues := strings.Split(enumElements[1], enumValueSplitter)

	return enumValues[rand.Intn(len(enumValues))]
}

func NewEnumProducer(tag string) Producer {
	return &EnumProducer{
		enumTag: tag,
	}
}
