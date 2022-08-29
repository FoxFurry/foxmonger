package main

import (
	"github.com/jaswdr/faker"
)

type FullNameProducer struct {
	fakerInstance faker.Faker
}

func (gen *FullNameProducer) Produce() string {
	return gen.fakerInstance.Person().Name()
}

func NewFullNameProducer() Producer {
	return &FullNameProducer{
		fakerInstance: faker.New(),
	}
}
