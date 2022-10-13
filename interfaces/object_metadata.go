package interfaces

import (
	"time"
)

type ObjectMetadata interface {
	Mode() ObjectMode
	ContentLength() *uint64
	ContentMD5() *string
	LastModified() *time.Time
	ETag() *string
}
