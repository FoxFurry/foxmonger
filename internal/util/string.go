package util

import (
	"strings"
)

var TagSplitter = ";"

func SplitTags(tagSet string) ([]string, error) {
	return strings.Split(tagSet, TagSplitter), nil
}
