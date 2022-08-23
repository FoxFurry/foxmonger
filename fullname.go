package main

import (
	"github.com/jaswdr/faker"
)

type FullNameGenerator struct{}

func (gen *FullNameGenerator) Generate() string {
	fakerInstance := faker.New()

	return fakerInstance.Person().Name()
}
