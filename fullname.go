package main

import (
	"github.com/jaswdr/faker"
)

type FullNameGenerator struct {
	fakerInstance faker.Faker
}

func (gen *FullNameGenerator) Generate() string {
	return gen.fakerInstance.Person().Name()
}

func NewFullNameGenerator() Generator {
	return &FullNameGenerator{
		fakerInstance: faker.New(),
	}
}
