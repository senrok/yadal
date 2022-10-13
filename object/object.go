package object

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	dalErrors "github.com/senrok/yadal/errors"
	"github.com/senrok/yadal/interfaces"
	"github.com/senrok/yadal/options"
	"github.com/senrok/yadal/utils"
	"io"
	"strings"
)

type Object struct {
	accessor interfaces.Accessor
	path     string
}

func NewObject(a interfaces.Accessor, p string) Object {
	return Object{
		accessor: a,
		path:     p,
	}
}

func (o *Object) ID() string {
	return fmt.Sprintf("%s/%s", o.accessor.Metadata().Root(), o.path)
}

func (o *Object) Path() string {
	return o.path
}

func (o *Object) Name() string {
	return utils.GetNameFromPath(o.path)
}

func (o *Object) Create(ctx context.Context) error {
	if strings.HasPrefix(o.path, "/") {
		return o.accessor.Create(ctx, o.path, options.CreateOptions{Mode: int8(interfaces.DIR)})
	}
	return o.accessor.Create(ctx, o.path, options.CreateOptions{Mode: int8(interfaces.FILE)})
}

func (o *Object) Read(ctx context.Context) (io.ReadCloser, error) {
	return o.RangeRead(ctx, nil)
}

func (o *Object) RangeRead(ctx context.Context, bytesRange options.RangeBounds) (io.ReadCloser, error) {
	opt := options.ReadOptions{}
	if bytesRange != nil {
		opt.Offset = bytesRange.Offset()
		opt.Size = bytesRange.Size()
	}
	return o.accessor.Read(ctx, o.path, opt)
}

var (
	ErrTryWrite2Dir = errors.New("try write bytes to a dir")
	ErrNotADir      = errors.New("not a directory")
)

func (o *Object) Write(ctx context.Context, byte []byte) error {
	if strings.HasSuffix(o.path, "/") {
		return ErrTryWrite2Dir
	}
	body := bytes.NewReader(byte)
	_, err := o.accessor.Write(ctx, o.path, options.WriteOptions{Size: uint64(len(byte))}, body)
	if err != nil {
		return err
	}
	return nil
}

func (o *Object) Delete(ctx context.Context) error {
	return o.accessor.Delete(ctx, o.path, options.DeleteOptions{})
}

func (o *Object) List(ctx context.Context) (interfaces.ObjectStream, error) {
	if !strings.HasPrefix(o.path, "/") {
		return nil, ErrNotADir
	}
	return o.accessor.List(ctx, o.path, options.ListOptions{})
}

func (o *Object) Metadata(ctx context.Context) (interfaces.ObjectMetadata, error) {
	return o.accessor.Stat(ctx, o.path, options.StatOptions{})
}

func (o *Object) IsExist(ctx context.Context) (bool, error) {
	_, err := o.accessor.Stat(ctx, o.path, options.StatOptions{})
	if err != nil {
		if errors.Is(err, dalErrors.ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
