// Package tag exposes Producer and Modifier interfaces describing row generation. Implementations of these interfaces
// will correspond to their tags
package tag

const (
	modifierSplitter  = ":"
	enumValueSplitter = ","
)

type Tag interface {
	Initialize(string) error
}

type Producer interface {
	Tag
	Produce() string
}

type Modifier interface {
	Tag
	Modify(string) string
}
