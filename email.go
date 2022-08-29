package main

import (
	"github.com/jaswdr/faker"
)

type EmailGenerator struct {
	fakerInstance faker.Faker
}

func (gen *EmailGenerator) Generate() string {
	return gen.fakerInstance.Internet().Email()
}

func NewEmailGenerator() Generator {
	return &EmailGenerator{
		fakerInstance: faker.New(),
	}
}
