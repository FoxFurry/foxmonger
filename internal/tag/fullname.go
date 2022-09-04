package tag

import (
	"github.com/jaswdr/faker"
)

type FullNameProducer struct {
	fakerInstance faker.Faker
}

func (gen *FullNameProducer) Produce() string {
	return gen.fakerInstance.Person().Name()
}

// Initialize implements Producer interface
func (gen *FullNameProducer) Initialize(string) error { return nil }

func NewFullNameProducer() Producer {
	return &FullNameProducer{
		fakerInstance: faker.New(),
	}
}
