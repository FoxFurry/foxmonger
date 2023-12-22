package modifiers

import (
	"github.com/FoxFurry/foxmonger/internal/tag"
)

type LimitModifier struct {
	limit int
}

func (mod *LimitModifier) Modify(val string) string {
	if len(val) > mod.limit {
		return val[0:mod.limit]
	}
	return val
}

func (mod *LimitModifier) Initialize(input string) error {
	//limitElements := strings.Split(input, tag.modifierSplitter)
	//
	//if len(limitElements) != 2 {
	//	return fmt.Errorf("tag does not contain lmit value: %s", input)
	//}
	//
	//limitValue, err := strconv.Atoi(limitElements[1])
	//if err != nil {
	//	return fmt.Errorf("could not convert limit value to integer: %s", input)
	//}
	//
	//mod.limit = limitValue
	return nil
}

func NewLimitModifier(tag string) (tag.Modifier, error) {
	lim := &LimitModifier{}
	if err := lim.Initialize(tag); err != nil {
		return nil, err
	}

	return lim, nil
}
