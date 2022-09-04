package tag

import (
	"math/rand"
)

type ForeignGenerator struct {
	foreignKeys []string
}

func (gen *ForeignGenerator) Produce() string {
	return gen.foreignKeys[rand.Intn(len(gen.foreignKeys))]
}
