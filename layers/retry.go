package layers

import (
	"context"
	"github.com/Rican7/retry"
	"github.com/Rican7/retry/strategy"
	"github.com/senrok/yadal/errors"
	"github.com/senrok/yadal/interfaces"
	"github.com/senrok/yadal/options"
	"io"
	"net/http"
)

type RetryOptions struct {
	Strategies []strategy.Strategy
}

type retryAccessor struct {
	inner interfaces.Accessor
	RetryOptions
}

func RetryWhen(err error, needRetry func(err error) bool) error {
	if needRetry(err) {
		return err
	}
	return nil
}

func IsErrInterrupted(err error) bool {
	return errors.Is(err, errors.ErrInterrupted)
}

func (r retryAccessor) Metadata() interfaces.Metadata {
	return r.inner.Metadata()
}

func (r retryAccessor) Create(ctx context.Context, path string, args options.CreateOptions) (innerErr error) {
	_ = retry.Retry(func(_ uint) error {
		innerErr = r.inner.Create(ctx, path, args)
		return RetryWhen(innerErr, IsErrInterrupted)
	}, r.Strategies...)
	return
}

func (r retryAccessor) Read(ctx context.Context, path string, args options.ReadOptions) (read io.ReadCloser, innerErr error) {
	_ = retry.Retry(func(_ uint) error {
		read, innerErr = r.inner.Read(ctx, path, args)
		return RetryWhen(innerErr, IsErrInterrupted)
	}, r.Strategies...)
	return
}

func (r retryAccessor) Write(ctx context.Context, path string, args options.WriteOptions, reader io.ReadSeeker) (size uint64, innerErr error) {
	_ = retry.Retry(func(_ uint) error {
		size, innerErr = r.inner.Write(ctx, path, args, reader)
		return RetryWhen(innerErr, IsErrInterrupted)
	}, r.Strategies...)
	return
}

func (r retryAccessor) Stat(ctx context.Context, path string, args options.StatOptions) (meta interfaces.ObjectMetadata, innerErr error) {
	_ = retry.Retry(func(_ uint) error {
		meta, innerErr = r.inner.Stat(ctx, path, args)
		return RetryWhen(innerErr, IsErrInterrupted)
	}, r.Strategies...)
	return
}

func (r retryAccessor) Delete(ctx context.Context, path string, args options.DeleteOptions) (innerErr error) {
	_ = retry.Retry(func(_ uint) error {
		innerErr = r.inner.Delete(ctx, path, args)
		return RetryWhen(innerErr, IsErrInterrupted)
	}, r.Strategies...)
	return
}

func (r retryAccessor) List(ctx context.Context, path string, args options.ListOptions) (stream interfaces.ObjectStream, innerErr error) {
	_ = retry.Retry(func(_ uint) error {
		stream, innerErr = r.inner.List(ctx, path, args)
		return RetryWhen(innerErr, IsErrInterrupted)
	}, r.Strategies...)
	return
}

func (r retryAccessor) PreSign(ctx context.Context, path string, args options.PreSignOptions) (req *http.Request, innerErr error) {
	_ = retry.Retry(func(_ uint) error {
		req, innerErr = r.inner.PreSign(ctx, path, args)
		return RetryWhen(innerErr, IsErrInterrupted)
	}, r.Strategies...)
	return
}

func (r retryAccessor) CreateMultipart(ctx context.Context, path string, args options.CreateMultipart) (uploadId string, innerErr error) {
	_ = retry.Retry(func(_ uint) error {
		uploadId, innerErr = r.inner.CreateMultipart(ctx, path, args)
		return RetryWhen(innerErr, IsErrInterrupted)
	}, r.Strategies...)
	return
}

func (r retryAccessor) WriteMultipart(ctx context.Context, path string, args options.WriteMultipart, reader io.ReadSeeker) (part interfaces.ObjectPart, innerErr error) {
	_ = retry.Retry(func(_ uint) error {
		part, innerErr = r.inner.WriteMultipart(ctx, path, args, reader)
		return RetryWhen(innerErr, IsErrInterrupted)
	}, r.Strategies...)
	return
}

func (r retryAccessor) CompleteMultipart(ctx context.Context, path string, args options.CompleteMultipart) (innerErr error) {
	_ = retry.Retry(func(_ uint) error {
		innerErr = r.inner.CompleteMultipart(ctx, path, args)
		return RetryWhen(innerErr, IsErrInterrupted)
	}, r.Strategies...)
	return
}

func (r retryAccessor) AbortMultipart(ctx context.Context, path string, args options.AbortMultipart) (innerErr error) {
	_ = retry.Retry(func(_ uint) error {
		innerErr = r.inner.AbortMultipart(ctx, path, args)
		return RetryWhen(innerErr, IsErrInterrupted)
	}, r.Strategies...)
	return
}

type RetryOption func(r *RetryOptions)

func SetStrategy(s ...strategy.Strategy) RetryOption {
	return func(r *RetryOptions) {
		r.Strategies = append(r.Strategies, s...)
	}
}

// NewRetryLayer returns a retry layer
func NewRetryLayer(opts ...RetryOption) interfaces.Layer {
	op := RetryOptions{}
	for _, opt := range opts {
		opt(&op)
	}
	return func(accessor interfaces.Accessor) interfaces.Accessor {
		return &retryAccessor{
			inner:        accessor,
			RetryOptions: op,
		}
	}
}
