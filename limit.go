package main

import (
	"fmt"
	"strconv"
	"strings"
)

type LimitModifier struct {
	limitTag string
}

func (mod *LimitModifier) Modify(val string) string {
	limitElements := strings.Split(mod.limitTag, modifierSplitter)

	if len(limitElements) != 2 { // Valid enum should contain 2 elements: enum keyword and enum values
		fmt.Printf("ignoring invalid limit tag: %s\n", mod.limitTag)
		return val
	}

	limitValue, err := strconv.Atoi(limitElements[1])
	if err != nil {
		fmt.Printf("could not extract limit value from tag: %s\n", mod.limitTag)
		return val
	}

	if len(val) > limitValue {
		return val[0:limitValue]
	}
	return val
}
