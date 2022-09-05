package util

import (
	"fmt"
	"strings"
)

var TagSplitter = ";"

type ResolvedForeign struct {
	TargetTable,
	TargetRow string
}

func SplitTags(tagSet string) ([]string, error) {
	return strings.Split(tagSet, TagSplitter), nil
}

func ResolveForeignKey(foreignTag string) (*ResolvedForeign, error) {
	foreignElements := strings.Split(foreignTag, TagSplitter)

	if len(foreignElements) != 2 {
		return nil, fmt.Errorf("tag does not contain foreign key description: %s", foreignTag)
	}

	var tableName, rowName string

	_, err := fmt.Sscanf(foreignElements[1], "%s(%s)", &tableName, &rowName)
	if err != nil {
		return nil, fmt.Errorf("failed parsing foreign key: %w", err)
	}

	return &ResolvedForeign{
		TargetTable: tableName,
		TargetRow:   rowName,
	}, nil
}
