package main

import (
	"github.com/jaswdr/faker"
)

type EmailGenerator struct{}

func (gen *EmailGenerator) Generate() string {
	fakerInstance := faker.New()

	return fakerInstance.Internet().Email()
}
