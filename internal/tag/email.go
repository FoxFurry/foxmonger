package tag

import (
	"github.com/jaswdr/faker"
)

type EmailProducer struct {
	fakerInstance *faker.Faker
}

func (gen *EmailProducer) Produce() string {
	return gen.fakerInstance.Internet().Email()
}

func (gen *EmailProducer) Initialize(string) error { return nil }

func NewEmailProducer(f *faker.Faker) Producer {
	return &EmailProducer{
		fakerInstance: f,
	}
}
