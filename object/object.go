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
		path:     utils.NormalizePath(p),
	}
}

func (o *Object) ID() string {
	return fmt.Sprintf("%s%s", o.accessor.Metadata().Root(), o.path)
}

func (o *Object) Path() string {
	return o.path
}

func (o *Object) Name() string {
	return utils.GetNameFromPath(o.path)
}

// Create it creates an empty object, like using the following linux commands:
//	- `touch path/to/file`
//	- `mkdir path/to/dir/`
//
// behaviors:
//
//	- create on existing dir will succeed.
//	- create on existing file will overwrite and truncate it.
//
// create a file:
// 	acc, _ := newS3Accessor()
//	op := NewOperatorFromAccessor(acc)
//	object := op.Object("test")
//	_ = object.Create(context.TODO())
//
// create a dir:
// 	acc, _ := newS3Accessor()
//	op := NewOperatorFromAccessor(acc)
//	object := op.Object("test-dir/")
//	_ = object.Create(context.TODO())
func (o *Object) Create(ctx context.Context) error {
	if strings.HasSuffix(o.path, "/") {
		return o.accessor.Create(ctx, o.path, options.CreateOptions{Mode: int8(interfaces.DIR)})
	}
	return o.accessor.Create(ctx, o.path, options.CreateOptions{Mode: int8(interfaces.FILE)})
}

// Read It returns a io.ReadCloser holds the whole object.
//
// read a file:
// 	acc, _ := newS3Accessor()
//	op := NewOperatorFromAccessor(acc)
//	object := op.Object("test")
//	_ = object.Write(context.TODO(), []byte("Hello,World!"))
//	reader, _ := object.Read(context.TODO())
//
//	_, _ = io.ReadAll(reader)
func (o *Object) Read(ctx context.Context) (io.ReadCloser, error) {
	return o.RangeRead(ctx, nil)
}

// RangeRead it returns a io.ReadCloser holds specified range of object .
//
// range read:
// 	acc, _ := newS3Accessor()
//	op := NewOperatorFromAccessor(acc)
//	object := op.Object("test")
//	_ = object.Write(context.TODO(), []byte("Hello,World!"))
//	reader, _ := object.RangeRead(context.TODO(), options.NewRangeBounds(options.Range(0, 11)))
//	// object.RangeRead(context.TODO(), options.NewRangeBounds(options.Start(2)))
//	// object.RangeRead(context.TODO(), options.NewRangeBounds(options.End(11)))
//
//	_, _ = io.ReadAll(reader)
func (o *Object) RangeRead(ctx context.Context, bytesRange options.RangeBounds) (io.ReadCloser, error) {
	if strings.HasSuffix(o.path, "/") {
		return nil, ErrIsADir
	}
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
	ErrIsADir       = errors.New("is a directory")
)

//  Write it writes bytes into object.
//
// write text:
// 	acc, _ := newS3Accessor()
//	op := NewOperatorFromAccessor(acc)
//	object := op.Object("test")
//	_ = object.Write(context.TODO(), []byte("Hello,World!"))
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

// Delete it deletes object.
//
// delete:
// 	acc, _ := newS3Accessor()
//	op := NewOperatorFromAccessor(acc)
//	object := op.Object("test")
//	_ = object.Delete(context.TODO())
func (o *Object) Delete(ctx context.Context) error {
	return o.accessor.Delete(ctx, o.path, options.DeleteOptions{})
}

// List it lists current directory object, returns a interfaces.ObjectStream.
//
// list:
// 	acc, _ := newS3Accessor()
//	op := NewOperatorFromAccessor(acc)
//	object := op.Object("dir/")
//	stream, _ := object.List(context.TODO())
//	for stream.HasNext() {
//		entry, _ := stream.Next(context.TODO())
//		if entry != nil {
//			fmt.Println(entry.Path())
//			fmt.Println(entry.Metadata().LastModified())
//			fmt.Println(entry.Metadata().ContentLength())
//		}
//	}
func (o *Object) List(ctx context.Context) (interfaces.ObjectStream, error) {
	if !strings.HasSuffix(o.path, "/") {
		return nil, ErrNotADir
	}
	return o.accessor.List(ctx, o.path, options.ListOptions{})
}

// Metadata it returns object's metadata, returns a interfaces.ObjectMetadata
//
// fetch metadata:
//	acc, _ := newS3Accessor()
//	op := NewOperatorFromAccessor(acc)
//	object := op.Object("test")
//	meta, err := object.Metadata(context.TODO())
//	if err == errors.ErrNotFound {
//		fmt.Println("not found")
//		return
//	}
//	fmt.Println(meta.LastModified())
//	fmt.Println(meta.ETag())
//	fmt.Println(meta.ContentLength())
//	fmt.Println(meta.ContentMD5())
//	fmt.Println(meta.Mode())
func (o *Object) Metadata(ctx context.Context) (interfaces.ObjectMetadata, error) {
	return o.accessor.Stat(ctx, o.path, options.StatOptions{})
}

// IsExist returns true if object exits
//
// example:
// 	acc, _ := newS3Accessor()
//	op := NewOperatorFromAccessor(acc)
//	object := op.Object("test")
//	fmt.Println(object.IsExist(context.TODO()))
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
