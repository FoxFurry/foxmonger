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
	enumPat  = regexp.MustCompile("enum:.+")
	limitPat = regexp.MustCompile("limit:\\d+")
)

type tagManager struct {
}

func (t *tagManager) resolveTag(tagString string) (*tagGenerator, error) {
	generator := tagGenerator{}

	tags := strings.Split(tagString, tagSplitter)

	for _, tag := range tags {
		switch {
		case tag == "auto": // Auto is used to keep history of autoincrement rows
			return nil, nil

		case tag == "fullname":
			generator.setGenerator(NewFullNameGenerator())

		case tag == "email":
			generator.setGenerator(NewEmailGenerator())

		case enumPat.MatchString(tag):
			generator.setGenerator(NewEnumGenerator(tag))

		case limitPat.MatchString(tag):
			generator.addModifier(NewLimitModifier(tag))

		default:
			return nil, fmt.Errorf("unsupported tag: %s", tag)
		}
	}

	return &generator, nil
}
