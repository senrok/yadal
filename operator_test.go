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
	"github.com/senrok/yadal/providers/s3"
	"go.uber.org/zap"
	"io"
	"math/rand"
	"os"
	"time"
)

var (
	DAL_BUCKET            string
	DAL_ENDPOINT          string
	DAL_ACCESS_KEY_ID     string
	DAL_SECRET_ACCESS_KEY string
)

func newS3Accessor() (interfaces.Accessor, error) {
	return s3.NewDriver(context.Background(), s3.Options{
		Bucket:    os.Getenv("Bucket"),
		Endpoint:  os.Getenv("Endpoint"),
		Root:      os.Getenv("Root"),
		Region:    os.Getenv("Region"),
		AccessKey: os.Getenv("AccessKey"),
		SecretKey: os.Getenv("SecretKey"),
	})
}

func ExampleOperator_Layer_logging() {
	acc, _ := newS3Accessor()
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
	acc, _ := newS3Accessor()
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
	acc, _ := newS3Accessor()
	op := NewOperatorFromAccessor(acc)
	object := op.Object("test")
	fmt.Println(object.IsExist(context.TODO()))
}

func ExampleOperator_Object_metadata() {
	acc, _ := newS3Accessor()
	op := NewOperatorFromAccessor(acc)
	object := op.Object("test")
	meta, err := object.Metadata(context.TODO())
	if err == errors.ErrNotFound {
		fmt.Println("not found")
		return
	}
	fmt.Println(meta.LastModified())
	fmt.Println(meta.ETag())
	fmt.Println(meta.ContentLength())
	fmt.Println(meta.ContentMD5())
	fmt.Println(meta.Mode())
}

func ExampleOperator_Object_list() {
	acc, _ := newS3Accessor()
	op := NewOperatorFromAccessor(acc)
	object := op.Object("dir/")
	stream, _ := object.List(context.TODO())
	for stream.HasNext() {
		entry, _ := stream.Next(context.TODO())
		if entry != nil {
			fmt.Println(entry.Path())
			fmt.Println(entry.Metadata().LastModified())
			fmt.Println(entry.Metadata().ContentLength())
		}
	}
}

func ExampleOperator_Object_delete() {
	acc, _ := newS3Accessor()
	op := NewOperatorFromAccessor(acc)
	object := op.Object("test")
	_ = object.Delete(context.TODO())
}

func ExampleOperator_Object_write() {
	acc, _ := newS3Accessor()
	op := NewOperatorFromAccessor(acc)
	object := op.Object("test")
	_ = object.Write(context.TODO(), []byte("Hello,World!"))
}

func ExampleOperator_Object_rangeRead() {
	acc, _ := newS3Accessor()
	op := NewOperatorFromAccessor(acc)
	object := op.Object("test")
	_ = object.Write(context.TODO(), []byte("Hello,World!"))
	reader, _ := object.RangeRead(context.TODO(), options.NewRangeBounds(options.Range(0, 11)))
	// object.RangeRead(context.TODO(), options.NewRangeBounds(options.Start(2)))
	// object.RangeRead(context.TODO(), options.NewRangeBounds(options.End(11)))

	_, _ = io.ReadAll(reader)
}

func ExampleOperator_Object_read() {
	acc, _ := newS3Accessor()
	op := NewOperatorFromAccessor(acc)
	object := op.Object("test")
	_ = object.Write(context.TODO(), []byte("Hello,World!"))
	reader, _ := object.Read(context.TODO())

	_, _ = io.ReadAll(reader)
}

func ExampleOperator_Object_create() {
	acc, _ := newS3Accessor()
	op := NewOperatorFromAccessor(acc)
	object := op.Object("test")
	_ = object.Create(context.TODO())
}

func ExampleOperator_Object() {
	acc, _ := newS3Accessor()
	op := NewOperatorFromAccessor(acc)
	object := op.Object("test")
	fmt.Println(object.ID())
	fmt.Println(object.Path())
	fmt.Println(object.Name())
}

func ExampleNewOperatorFromAccessor() {
	acc, _ := newS3Accessor()
	op := NewOperatorFromAccessor(acc)
	// Create object handler
	o := op.Object("test_file")

	// Write data
	if err := o.Write(context.Background(), []byte("Hello,World!")); err != nil {
		return
	}

	// Read data
	bs, _ := o.Read(context.Background())
	fmt.Println(bs)
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

	// Read Dir
	ds := op.Object("test-dir/")
	iter, _ := ds.List(context.Background())
	for iter.HasNext() {
		entry, _ := iter.Next(context.Background())
		entry.Path()
		entry.Metadata()
	}
}
