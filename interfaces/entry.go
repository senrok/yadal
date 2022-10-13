package interfaces

import (
	"sync"
)

type entry struct {
	mu       *sync.Mutex
	accessor *Accessor
	path     string
	metadata ObjectMetadata
	complete bool
}

type Entry interface {
	Accessor() Accessor
	Path() string
	Metadata() ObjectMetadata
	IsComplete() bool
}
