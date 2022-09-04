package tag

import (
	"math/rand"
)

type ForeignProducer struct {
	foreignKeys []string
}

func (gen *ForeignProducer) Produce() string {
	return gen.foreignKeys[rand.Intn(len(gen.foreignKeys))]
}

func (gen *ForeignProducer) Initialize(string) error { return nil }

func NewForeignProducer(keys []string) Producer {
	return &ForeignProducer{
		foreignKeys: keys,
	}
}
