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

func (t *tagManager) resolveTag(tagString string) *tagGenerator {
	generator := tagGenerator{}

	tags := strings.Split(tagString, tagSplitter)

	for _, tag := range tags {
		switch {
		case tag == "auto":
			return nil

		case tag == "fullname":
			generator.setGenerator(&FullNameGenerator{})

		case tag == "email":
			generator.setGenerator(&EmailGenerator{})

		case enumPat.MatchString(tag):
			generator.setGenerator(&EnumGenerator{tag})

		case limitPat.MatchString(tag):
			generator.addModifier(&LimitModifier{tag})

		default:
			fmt.Printf("unsupported tag: %s\n", tag)
			return nil
		}
	}

	return &generator
}
