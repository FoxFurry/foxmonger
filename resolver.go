package main

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	modifierSplitter  = ":"
	enumValueSplitter = ","
)

var (
	enumPat    = regexp.MustCompile("enum:.+")
	limitPat   = regexp.MustCompile("limit:\\d+")
	foreignPat = regexp.MustCompile("foreign:.+")
)

type tagManager struct{}

func (t *tagManager) resolveTag(tagString string) (*Generator, error) {
	generator := NewEmptyGenerator()

	tags := strings.Split(tagString, tagSplitter)

	for _, tag := range tags {
		switch {

		// Exact tags
		case tag == "auto": // Auto is used to keep history of autoincrement rows
			return nil, nil

		case tag == "fullname":
			generator.setProducer(NewFullNameProducer())

		case tag == "email":
			generator.setProducer(NewEmailProducer())

		// Pattern tags
		case enumPat.MatchString(tag):
			generator.setProducer(NewEnumProducer(tag))

		case limitPat.MatchString(tag):
			generator.addModifier(NewLimitModifier(tag))

		case foreignPat.MatchString(tag):
			return nil, fmt.Errorf("foreign key tag is not yet supported")

		default:
			return nil, fmt.Errorf("unknown tag: %s", tag)
		}
	}

	return &generator, nil
}
