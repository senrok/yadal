package object

import "github.com/senrok/yadal/interfaces"

type Entry struct {
	accessor interfaces.Accessor
	path     string
	metadata interfaces.ObjectMetadata
	complete bool
}

func (e Entry) Accessor() interfaces.Accessor {
	return e.accessor
}

func (e Entry) Path() string {
	return e.path
}

func (e Entry) Metadata() interfaces.ObjectMetadata {
	return e.metadata
}

func (e Entry) IsComplete() bool {
	return e.complete
}

func NewEntry(accessor interfaces.Accessor, path string, metadata interfaces.ObjectMetadata, complete bool) *Entry {
	return &Entry{
		accessor: accessor,
		path:     path,
		metadata: metadata,
		complete: complete,
	}
}
