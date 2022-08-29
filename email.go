package main

import (
	"github.com/jaswdr/faker"
)

type EmailProducer struct {
	fakerInstance faker.Faker
}

func (gen *EmailProducer) Produce() string {
	return gen.fakerInstance.Internet().Email()
}

func NewEmailProducer() Producer {
	return &EmailProducer{
		fakerInstance: faker.New(),
	}
}
