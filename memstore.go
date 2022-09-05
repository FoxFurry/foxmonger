package foxmonger

import (
	"fmt"
)

type memstore struct {
	storage map[string][]string
}

func (m *memstore) AddValue(row, table, value string) {
	RowTableKey := MergeRowTable(row, table)
	m.storage[RowTableKey] = append(m.storage[RowTableKey], value)
}

func (m *memstore) GetValues(row, table string) []string {
	RowTableKey := MergeRowTable(row, table)
	return m.storage[RowTableKey]
}

func MergeRowTable(row, table string) string {
	return fmt.Sprintf("%s:%s", row, table)
}
