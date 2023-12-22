package producers

import (
	"github.com/FoxFurry/foxmonger/internal/tag"
	"github.com/jaswdr/faker"
)

type EmailProducer struct {
	fakerInstance *faker.Faker
}

func (gen *EmailProducer) Produce() string {
	return gen.fakerInstance.Internet().Email()
}

func (gen *EmailProducer) Initialize(string) error { return nil }

func NewEmailProducer(f *faker.Faker) tag.Producer {
	return &EmailProducer{
		fakerInstance: f,
	}
}
