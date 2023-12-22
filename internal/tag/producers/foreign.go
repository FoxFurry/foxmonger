package producers

import (
	"math/rand"
	"strings"

	"github.com/FoxFurry/foxmonger/internal/tag"
)

type ForeignProducer struct {
	foreignValuePool []string
}

func (gen *ForeignProducer) Produce() string {
	return gen.foreignValuePool[rand.Intn(len(gen.foreignValuePool))]
}

func (gen *ForeignProducer) Initialize(foreigns string) error {
	gen.foreignValuePool = strings.Split(foreigns, ",")

	return nil
}

func NewForeignProducer(foreigns []string) tag.Producer {
	return &ForeignProducer{
		foreignValuePool: foreigns,
	}
}
