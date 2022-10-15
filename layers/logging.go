package layers

import (
	"context"
	"github.com/senrok/yadal/interfaces"
	"github.com/senrok/yadal/options"
	"github.com/senrok/yadal/utils"
	"io"
	"net/http"
)

type adapter struct {
	info  func(args ...interface{})
	infof func(template string, args ...interface{})
}

func (a adapter) Info(args ...interface{}) {
	a.info(args...)
}

func (a adapter) Infof(template string, args ...interface{}) {
	a.infof(template, args...)
}

func NewLoggerAdapter(info func(args ...interface{}), infof func(template string, args ...interface{})) Logger {
	return &adapter{
		info:  info,
		infof: infof,
	}
}

type Logger interface {
	Info(args ...interface{})
	Infof(template string, args ...interface{})
}

type LoggingOptions struct {
	Logger
}

type LoggingOption func(r *LoggingOptions)

type loggingAccessor struct {
	inner interfaces.Accessor
	Logger
}

func (l loggingAccessor) innerProvider() string {
	return l.inner.Metadata().Provider().String()
}

func (l loggingAccessor) Metadata() interfaces.Metadata {
	l.Infof("dal::service service=%s operation=%s -> starting", interfaces.MetadataOP, l.innerProvider())
	meta := l.inner.Metadata()
	l.Infof("dal::service service=%s operation=%s -> finished: %s", interfaces.MetadataOP, l.innerProvider(), meta)
	return meta
}

func (l loggingAccessor) Create(ctx context.Context, path string, args options.CreateOptions) error {
	l.Infof("dal::service service=%s operation=%s -> starting", interfaces.CreateOp, l.innerProvider())
	err := l.inner.Create(ctx, path, args)
	l.Infof("dal::service service=%s operation=%s -> finished", interfaces.CreateOp, l.innerProvider())
	if err != nil {
		l.Infof("dal::service service=%s operation=%s -> error: %s", interfaces.CreateOp, l.innerProvider(), err)
	}
	return err
}

func (l loggingAccessor) Read(ctx context.Context, path string, args options.ReadOptions) (io.ReadCloser, error) {
	l.Infof("dal::service service=%s operation=%s size=%s offset=%s -> starting", interfaces.ReadOp, l.innerProvider(), utils.FmtPtrUint64(args.Size), utils.FmtPtrUint64(args.Offset))
	reader, err := l.inner.Read(ctx, path, args)
	l.Infof("dal::service service=%s operation=%s -> finished", interfaces.ReadOp, l.innerProvider())
	if err != nil {
		l.Infof("dal::service service=%s operation=%s -> error: %s", interfaces.ReadOp, l.innerProvider(), err)
	}
	return reader, err
}

func (l loggingAccessor) Write(ctx context.Context, path string, args options.WriteOptions, reader io.ReadSeeker) (uint64, error) {
	l.Infof("dal::service service=%s operation=%s size=%s -> starting", interfaces.WriteOp, l.innerProvider(), args.Size)
	size, err := l.inner.Write(ctx, path, args, reader)
	l.Infof("dal::service service=%s operation=%s -> finished", interfaces.WriteOp, l.innerProvider())
	if err != nil {
		l.Infof("dal::service service=%s operation=%s -> error: %s", interfaces.WriteOp, l.innerProvider(), err)
	}
	return size, err
}

func (l loggingAccessor) Stat(ctx context.Context, path string, args options.StatOptions) (interfaces.ObjectMetadata, error) {
	l.Infof("dal::service service=%s operation=%s -> starting", interfaces.StatOp, l.innerProvider())
	meta, err := l.inner.Stat(ctx, path, args)
	l.Infof("dal::service service=%s operation=%s -> finished", interfaces.StatOp, l.innerProvider())
	if err != nil {
		l.Infof("dal::service service=%s operation=%s -> error: %s", interfaces.StatOp, l.innerProvider(), err)
	}
	return meta, err
}

func (l loggingAccessor) Delete(ctx context.Context, path string, args options.DeleteOptions) error {
	l.Infof("dal::service service=%s operation=%s -> starting", interfaces.DeleteOp, l.innerProvider())
	err := l.inner.Delete(ctx, path, args)
	l.Infof("dal::service service=%s operation=%s -> finished", interfaces.DeleteOp, l.innerProvider())
	if err != nil {
		l.Infof("dal::service service=%s operation=%s -> error: %s", interfaces.DeleteOp, l.innerProvider(), err)
	}
	return err
}

func (l loggingAccessor) List(ctx context.Context, path string, args options.ListOptions) (interfaces.ObjectStream, error) {
	l.Infof("dal::service service=%s operation=%s -> starting", interfaces.ListOp, l.innerProvider())
	stream, err := l.inner.List(ctx, path, args)
	l.Infof("dal::service service=%s operation=%s -> finished", interfaces.ListOp, l.innerProvider())
	if err != nil {
		l.Infof("dal::service service=%s operation=%s -> error: %s", interfaces.ListOp, l.innerProvider(), err)
	}
	return stream, err
}

func (l loggingAccessor) PreSign(ctx context.Context, path string, args options.PreSignOptions) (*http.Request, error) {
	l.Infof("dal::service service=%s operation=%s -> starting", interfaces.PreSignOp, l.innerProvider())
	req, err := l.inner.PreSign(ctx, path, args)
	l.Infof("dal::service service=%s operation=%s -> finished", interfaces.PreSignOp, l.innerProvider())
	if err != nil {
		l.Infof("dal::service service=%s operation=%s -> error: %s", interfaces.PreSignOp, l.innerProvider(), err)
	}
	return req, err
}

func (l loggingAccessor) CreateMultipart(ctx context.Context, path string, args options.CreateMultipart) (string, error) {
	l.Infof("dal::service service=%s operation=%s -> starting", interfaces.CreateMultipartOp, l.innerProvider())
	uploadId, err := l.inner.CreateMultipart(ctx, path, args)
	l.Infof("dal::service service=%s operation=%s -> finished", interfaces.CreateMultipartOp, l.innerProvider())
	if err != nil {
		l.Infof("dal::service service=%s operation=%s -> error: %s", interfaces.CreateMultipartOp, l.innerProvider(), err)
	}
	return uploadId, err
}

func (l loggingAccessor) WriteMultipart(ctx context.Context, path string, args options.WriteMultipart, reader io.ReadSeeker) (interfaces.ObjectPart, error) {
	l.Infof("dal::service service=%s operation=%s -> starting", interfaces.WriteMultipartOp, l.innerProvider())
	part, err := l.inner.WriteMultipart(ctx, path, args, reader)
	l.Infof("dal::service service=%s operation=%s -> finished", interfaces.WriteMultipartOp, l.innerProvider())
	if err != nil {
		l.Infof("dal::service service=%s operation=%s -> error: %s", interfaces.WriteMultipartOp, l.innerProvider(), err)
	}
	return part, err
}

func (l loggingAccessor) CompleteMultipart(ctx context.Context, path string, args options.CompleteMultipart) error {
	l.Infof("dal::service service=%s operation=%s -> starting", interfaces.WriteMultipartOp, l.innerProvider())
	err := l.inner.CompleteMultipart(ctx, path, args)
	l.Infof("dal::service service=%s operation=%s -> finished", interfaces.WriteMultipartOp, l.innerProvider())
	if err != nil {
		l.Infof("dal::service service=%s operation=%s -> error: %s", interfaces.WriteMultipartOp, l.innerProvider(), err)
	}
	return err
}

func (l loggingAccessor) AbortMultipart(ctx context.Context, path string, args options.AbortMultipart) error {
	l.Infof("dal::service service=%s operation=%s -> starting", interfaces.AbortMultipartOp, l.innerProvider())
	err := l.inner.AbortMultipart(ctx, path, args)
	l.Infof("dal::service service=%s operation=%s -> finished", interfaces.AbortMultipartOp, l.innerProvider())
	if err != nil {
		l.Infof("dal::service service=%s operation=%s -> error: %s", interfaces.AbortMultipartOp, l.innerProvider(), err)
	}
	return err
}

func SetLogger(logger Logger) LoggingOption {
	return func(r *LoggingOptions) {
		r.Logger = logger
	}
}

// NewLoggingLayer returns logging layer
func NewLoggingLayer(opts ...LoggingOption) interfaces.Layer {
	op := LoggingOptions{}
	for _, opt := range opts {
		opt(&op)
	}
	return func(accessor interfaces.Accessor) interfaces.Accessor {
		return &loggingAccessor{inner: accessor, Logger: op.Logger}
	}
}
