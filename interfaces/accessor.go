package interfaces

import (
	"context"
	"github.com/senrok/yadal/options"
	"io"
	"net/http"
	"strings"
)

type Accessor interface {
	Metadata() Metadata

	// Create
	//
	// # Behavior
	// 	 - Input path MUST match with object.ObjectMode, WITHOUT checking ObjectMode.
	// 	 - Creating on existing dir SHOULD succeed.
	// 	 - Creating on existing file SHOULD overwrite and truncate.
	Create(ctx context.Context, path string, args options.CreateOptions) error

	// Read returns [io.ReadCloser] if operation succeeded
	//
	// # Behavior
	// 	 - Input path MUST match with object.ObjectMode, WITHOUT checking ObjectMode.
	Read(ctx context.Context, path string, args options.ReadOptions) (io.ReadCloser, error)

	// Write returns written size if operation succeeded
	//
	// # Behavior
	// 	 - Input path MUST be file path, WITHOUT checking ObjectMode.
	Write(ctx context.Context, path string, args options.WriteOptions, reader io.ReadSeeker) (uint64, error)

	// Stat
	//
	//	- Stat an empty path means stat provider's root path.
	// 	- Stat a path ends-with "/" means stating a dir.
	// 	- `mode` and `content_length` must be set.
	Stat(ctx context.Context, path string, args options.StatOptions) (ObjectMetadata, error)

	// Delete
	//
	// # Behavior
	//
	// 	- it is an idempotent operation, it's safe to call `Delete` on the same path multiple times.
	// 	- it SHOULD return nil if the path is deleted successfully or not exist.
	Delete(ctx context.Context, path string, args options.DeleteOptions) error

	// List
	//
	// # Behavior
	//
	//  - Input path MUST be dir path, DON'T NEED to check ObjectMode.
	//  - List non-exist dir should return Empty.
	List(ctx context.Context, path string, args options.ListOptions) (ObjectStream, error)

	// PreSign
	//
	// # Behavior
	//
	//	- Requires capability: `PreSign`
	// 	- This API is optional, throws errors.ErrUnsupportedMethod if not supported.
	PreSign(ctx context.Context, path string, args options.PreSignOptions) (*http.Request, error)

	// CreateMultipart
	//
	// # Behavior
	//
	//  - Requires capability: `Multipart`
	// 	- This op returns a `upload_id` which is required to for following APIs.
	CreateMultipart(ctx context.Context, path string, args options.CreateMultipart) (string, error)

	// WriteMultipart
	//
	// # Behavior
	//
	//  - Requires capability: `Multipart`
	WriteMultipart(ctx context.Context, path string, args options.WriteMultipart, reader io.ReadSeeker) (ObjectPart, error)

	// CompleteMultipart
	// # Behavior
	//
	//  - Requires capability: `Multipart`
	CompleteMultipart(ctx context.Context, path string, args options.CompleteMultipart) error

	// AbortMultipart
	// # Behavior
	//
	//  - Requires capability: `Multipart`
	AbortMultipart(ctx context.Context, path string, args options.AbortMultipart) error
}

type Capability uint8

func (c Capability) Has(capabilities ...Capability) bool {
	for _, capability := range capabilities {
		if c&capability == 0 {
			return false
		}
	}
	return true
}

func (c Capability) String() string {
	var result []string
	for _, capability := range cRange {
		if c.Has(capability) {
			result = append(result, capability.CapString())
		}
	}
	return strings.Join(result, "|")
}

var (
	cRange = []Capability{Read, Write, List, PreSign, Multipart, Blocking}
)

func (c Capability) CapString() string {
	switch c {
	case Read:
		return "Read"
	case Write:
		return "Write"
	case List:
		return "List"
	case PreSign:
		return "PreSign"
	case Multipart:
		return "Multipart"
	case Blocking:
		return "Blocking"
	default:
		return "Unknown"
	}
}

const (
	// Read `read` and `stat`
	Read Capability = 1 << iota

	// Write `write` and `delete`
	Write

	// List `list`
	List

	// PreSign `preSign`
	PreSign

	// Multipart `multipart`
	Multipart

	// Blocking `blocking`
	Blocking
)

type Metadata interface {
	Provider() Provider
	Root() string
	Name() string
	Capability() Capability
	String() string
}
