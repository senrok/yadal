package layers

import (
	"context"
	"github.com/senrok/yadal/interfaces"
	"github.com/senrok/yadal/options"
	"io"
	"net/http"
)

type Base struct {
}

type Ctx struct {
	Ctx                      *context.Context
	Path                     string
	Method                   string
	InputBuffer              *io.ReadCloser
	OutputBuffer             *io.ReadCloser
	CreateOptions            *options.CreateOptions
	ReadOptions              *options.ReadOptions
	WriteOptions             *options.WriteOptions
	StatOptions              *options.StatOptions
	DeleteOptions            *options.DeleteOptions
	ListOptions              *options.ListOptions
	PreSignOptions           *options.PreSignOptions
	CreateMultipartOptions   *options.CreateMultipart
	WriteMultipartOptions    *options.WriteMultipart
	CompleteMultipartOptions *options.CompleteMultipart
	AbortMultipartOptions    *options.AbortMultipart
	ObjectMetadata           interfaces.ObjectMetadata
	ObjectStream             interfaces.ObjectStream
	HttpRequest              *http.Request
	Err                      error
}

type BaseOptions struct {
	BeforeMetadata func(ctx *Ctx) interfaces.Metadata
	AfterMetadata  func(ctx *Ctx) interfaces.Metadata

	BeforeCreate func(ctx *Ctx) error
	AfterCreate  func(ctx *Ctx) error

	BeforeWrite func(ctx *Ctx) (uint64, error)
	AfterWrite  func(ctx *Ctx) (uint64, error)

	BeforeStat func(ctx *Ctx) (interfaces.ObjectMetadata, error)
	AfterStat  func(ctx *Ctx) (interfaces.ObjectMetadata, error)

	BeforeDelete func(ctx *Ctx) error
	AfterDelete  func(ctx *Ctx) error

	BeforeList func(ctx *Ctx) (interfaces.ObjectStream, error)
	AfterList  func(ctx *Ctx) (interfaces.ObjectStream, error)
}

type BaseOption func(o *BaseOptions)

func NewBaseLayer(opts ...BaseOption) interfaces.Layer {
	return func(accessor interfaces.Accessor) interfaces.Accessor {
		return nil
	}
}
