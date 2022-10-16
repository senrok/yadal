package yadal

import (
	"context"
	"fmt"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/jitter"
	"github.com/Rican7/retry/strategy"
	"github.com/senrok/yadal/errors"
	"github.com/senrok/yadal/interfaces"
	"github.com/senrok/yadal/layers"
	"github.com/senrok/yadal/options"
	"github.com/senrok/yadal/providers/fs"
	"go.uber.org/zap"
	"io"
	"math/rand"
	"time"
)

var (
	DAL_BUCKET            string
	DAL_ENDPOINT          string
	DAL_ACCESS_KEY_ID     string
	DAL_SECRET_ACCESS_KEY string
)

func newFsAccessor() (interfaces.Accessor, error) {
	return fs.NewDriver(fs.Options{Root: "tmp/"}), nil
}

func ExampleOperator_Layer_logging() {
	acc, _ := newFsAccessor()
	op := NewOperatorFromAccessor(acc)

	// logger
	logger, _ := zap.NewProduction()
	s := logger.Sugar()

	// logging layer
	loggingLayer := layers.NewLoggingLayer(
		layers.SetLogger(
			layers.NewLoggerAdapter(s.Info, s.Infof),
		),
	)

	op.Layer(loggingLayer)
}

func ExampleOperator_Layer_retry() {
	acc, _ := newFsAccessor()
	op := NewOperatorFromAccessor(acc)

	seed := time.Now().UnixNano()
	random := rand.New(rand.NewSource(seed))

	// retry layer
	// see strategies: https://github.com/Rican7/retry
	retryLayer := layers.NewRetryLayer(
		layers.SetStrategy(
			strategy.Limit(5),
			strategy.BackoffWithJitter(
				backoff.BinaryExponential(10*time.Millisecond),
				jitter.Deviation(random, 0.5),
			),
		),
	)

	op.Layer(retryLayer)
}

func ExampleOperator_Object_isExist() {
	acc, _ := newFsAccessor()
	op := NewOperatorFromAccessor(acc)
	object := op.Object("test")
	exist, _ := object.IsExist(context.TODO())
	fmt.Println(exist)

	// Output: false
}

func ExampleOperator_Object_metadata() {
	acc, _ := newFsAccessor()
	op := NewOperatorFromAccessor(acc)
	object := op.Object("test")
	_ = object.Create(context.TODO())
	meta, err := object.Metadata(context.TODO())
	if err == errors.ErrNotFound {
		fmt.Println("not found")
		return
	}
	//fmt.Println(meta.LastModified())
	fmt.Println(*meta.ETag())
	fmt.Println(*meta.ContentLength())
	fmt.Println(*meta.ContentMD5())
	fmt.Println(meta.Mode())
	fmt.Println(object.Path())

	// Output:0
	//
	// file
	// test
}

func ExampleOperator_Object_list() {
	acc, _ := newFsAccessor()
	op := NewOperatorFromAccessor(acc)
	o := op.Object("dir/test")
	_ = o.Create(context.TODO())
	object := op.Object("dir/")
	stream, _ := object.List(context.TODO())
	for stream.HasNext() {
		entry, _ := stream.Next(context.TODO())
		if entry != nil {
			fmt.Println(entry.Path())
			//fmt.Println(entry.Metadata().LastModified())
			fmt.Println(*entry.Metadata().ContentLength())
		}
	}

	// Output: dir/test
	// 0
}

func ExampleOperator_Object_delete() {
	acc, _ := newFsAccessor()
	op := NewOperatorFromAccessor(acc)
	object := op.Object("test")
	_ = object.Delete(context.TODO())
}

func ExampleOperator_Object_write() {
	acc, _ := newFsAccessor()
	op := NewOperatorFromAccessor(acc)
	object := op.Object("test")
	_ = object.Write(context.TODO(), []byte("Hello,World!"))
	reader, _ := object.Read(context.TODO())
	bytes, _ := io.ReadAll(reader)
	fmt.Println(string(bytes))

	//Output: Hello,World!
}

func ExampleOperator_Object_rangeRead() {
	acc, _ := newFsAccessor()
	op := NewOperatorFromAccessor(acc)
	object := op.Object("test")
	_ = object.Write(context.TODO(), []byte("Hello,World!"))
	reader, _ := object.RangeRead(context.TODO(), options.NewRangeBounds(options.Range(3, 8)))
	// object.RangeRead(context.TODO(), options.NewRangeBounds(options.Start(2)))
	// object.RangeRead(context.TODO(), options.NewRangeBounds(options.End(11)))

	bytes, _ := io.ReadAll(reader)
	fmt.Println(string(bytes))

	//Output: lo,Wo
}

func ExampleOperator_Object_read() {
	acc, _ := newFsAccessor()
	op := NewOperatorFromAccessor(acc)
	object := op.Object("test")
	_ = object.Write(context.TODO(), []byte("Hello,World!"))
	reader, _ := object.Read(context.TODO())

	bytes, _ := io.ReadAll(reader)
	fmt.Println(string(bytes))

	// Output: Hello,World!
}

func ExampleOperator_Object_create() {
	acc, _ := newFsAccessor()
	op := NewOperatorFromAccessor(acc)
	object := op.Object("test")
	_ = object.Create(context.TODO())
	fmt.Println(object.ID())
	fmt.Println(object.Path())

	// Output: /tmp/test
	// test
}

func ExampleOperator_Object() {
	acc, _ := newFsAccessor()
	op := NewOperatorFromAccessor(acc)
	object := op.Object("dir/test")
	fmt.Println(object.ID())
	fmt.Println(object.Path())
	fmt.Println(object.Name())

	// Output: /tmp/dir/test
	// dir/test
	// test
}

func ExampleNewOperatorFromAccessor() {
	acc, _ := newFsAccessor()
	op := NewOperatorFromAccessor(acc)
	// Create object handler
	o := op.Object("test_file")

	// Write data
	if err := o.Write(context.Background(), []byte("Hello,World!")); err != nil {
		return
	}

	// Read data
	bs, _ := o.Read(context.Background())
	bytes, _ := io.ReadAll(bs)
	fmt.Println(string(bytes))
	name := o.Name()
	fmt.Println(name)
	path := o.Path()
	fmt.Println(path)
	meta, _ := o.Metadata(context.Background())

	_ = meta.ETag()
	_ = meta.ContentLength()
	_ = meta.ContentMD5()
	_ = meta.Mode()

	// Delete
	_ = o.Delete(context.Background())

	dirFile := op.Object("test-dir/test")
	dirFile.Create(context.TODO())

	// Read Dir
	ds := op.Object("test-dir/")
	iter, _ := ds.List(context.Background())
	for iter.HasNext() {
		entry, _ := iter.Next(context.Background())
		if entry != nil {
			fmt.Println(entry.Path())
			fmt.Println(*entry.Metadata().ContentLength())
		}
	}

	// Output: Hello,World!
	// test_file
	// test_file
	// test-dir/test
	// 0
}
