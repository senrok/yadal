<p align="center"><img src="./assets/sqaure-logo.png" width="370"></p>
<p align="center">
<b>A project of SENROK Open Source</b>
</p>


# YaDAL 

<p>
<a href="https://goreportcard.com/report/github.com/senrok/yadal">
<img src="https://goreportcard.com/badge/github.com/senrok/yadal">
</a>
<a href="https://godoc.org/github.com/senrok/yadal">
<img src="https://godoc.org/github.com/senrok/yadal?status.svg" alt="GoDoc">
</a>
<a href="https://github.com/senrok/yadal/actions/workflows/service_test_s3.yml">
<img src="https://github.com/senrok/yadal/actions/workflows/service_test_s3.yml/badge.svg"/>
</a>
<a href="https://github.com/senrok/yadal/actions/workflows/service_test_fs.yml">
<img src="https://github.com/senrok/yadal/actions/workflows/service_test_fs.yml/badge.svg"/>
</a>
</p>


**Y**et **A**nother **D**ata **A**ccess **L**ayer: Access data freely, efficiently, without the tears 😢

inspired by [Databend's OpenDAL](https://github.com/datafuselabs/opendal)

## Table of contents

- [Features](#features)
- [Installation](#installation)
- [Get started](#get-started)
- [Documentation](#documentation)
  - [Object](#object)
    - [Handler](#handler)
    - [IsExist](#isexist)
    - [Metadata](#metadata)
    - [Create a dir or a file](#create-a-dir-or-a-object)
    - [Read](#read)
    - [Range Read](#range-read)
    - [Write](#write)
    - [Delete](#delete)
    - [List current directory](#list-current-directory)
  - [Layers](#layers)
    - [Retry](#retry)
    - [Logging](#logging)
  
- [License](#license)

## Features

**Freely**
- [x] Access different storage services in the same way
- [ ] Behavior tests for all services
  - [x] S3 and S3 compatible services
  - [x] fs: POSIX compatible filesystem

**Without the tears 😢**
- [x] Powerful Layer Middlewares
  - [x] Auto Retry (Backoff)
  - [x] Logging Layer
  - [ ] Tracing Layer
  - [ ] Metrics Layer
- [ ] Compress/Decompress 
- [ ] Service-side encryption

**Efficiently**
- Zero cost: mapping to underlying API calls directly
- Auto metadata reuse: avoid extra metadata calls


## Installation

```bash
go get -u github.com/senrok/yadal
```

## Get started

```go
import (
	"context"
	"fmt"
	"github.com/senrok/yadal"
	"github.com/senrok/yadal/providers/s3"
	"os"
)

func main() {
	acc, _ := s3.NewDriver(context.Background(), s3.Options{
		Bucket:    os.Getenv("Bucket"),
		Endpoint:  os.Getenv("Endpoint"),
		Root:      os.Getenv("Root"),
		Region:    os.Getenv("Region"),
		AccessKey: os.Getenv("AccessKey"),
		SecretKey: os.Getenv("SecretKey"),
	})
	op := yadal.NewOperatorFromAccessor(acc)
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

```

See the [Documentation](https://godoc.org/github.com/senrok/yadal) or explore more [examples](examples)

## Documentation

<a href="https://godoc.org/github.com/senrok/yadal">
<img src="https://godoc.org/github.com/senrok/yadal?status.svg" alt="GoDoc">
</a>

### Object
#### Handler
```go
func ExampleOperator_Object() {
	acc, _ := newS3Accessor()
	op := NewOperatorFromAccessor(acc)
	// it returns a object.Object handler 
	object := op.Object("test")
	fmt.Println(object.ID())
	fmt.Println(object.Path())
	fmt.Println(object.Name())
}
```
#### Create a dir or a object
It creates an empty object, like using the following linux commands:
- `touch path/to/file`
- `mkdir path/to/dir/`

The behaviors: 
- create on existing dir will succeed.
- create on existing file will overwrite and truncate it.

a dir:
```go
func ExampleOperator_Object_create() {
	acc, _ := newS3Accessor()
	op := NewOperatorFromAccessor(acc)
	object := op.Object("test/")
	_ = object.Create(context.TODO())
}
```

a object:
```go
func ExampleOperator_Object_create() {
	acc, _ := newS3Accessor()
	op := NewOperatorFromAccessor(acc)
	object := op.Object("test")
	_ = object.Create(context.TODO())
}
```

#### IsExist

```go
func ExampleOperator_Object_isExist() {
	acc, _ := newS3Accessor()
	op := NewOperatorFromAccessor(acc)
	object := op.Object("test")
	fmt.Println(object.IsExist(context.TODO()))
}
```

#### Metadata

```go
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
```

#### Read

It returns a io.ReadCloser holds the whole object.

```go
func ExampleOperator_Object_read() {
	acc, _ := newS3Accessor()
	op := NewOperatorFromAccessor(acc)
	object := op.Object("test")
	_ = object.Write(context.TODO(), []byte("Hello,World!"))
	reader, _ := object.Read(context.TODO())

	_, _ = io.ReadAll(reader)
}
```


#### Range Read

It returns a io.ReadCloser holds specified range of object .

```go
func ExampleOperator_Object_rangeRead() {
	acc, _ := newS3Accessor()
	op := NewOperatorFromAccessor(acc)
	object := op.Object("test")
	_ = object.Write(context.TODO(), []byte("Hello,World!"))
	reader, _ := object.RangeRead(context.TODO(), options.NewRangeBounds(options.Range(0, 11)))
	_, _ = object.RangeRead(context.TODO(), options.NewRangeBounds(options.Start(2)))
	_, _ = object.RangeRead(context.TODO(), options.NewRangeBounds(options.End(11)))

	_, _ = io.ReadAll(reader)
}
```

#### Write

It writes bytes into object.

```go
func ExampleOperator_Object_write() {
	acc, _ := newS3Accessor()
	op := NewOperatorFromAccessor(acc)
	object := op.Object("test")
	_ = object.Write(context.TODO(), []byte("Hello,World!"))
}
```

#### Delete

It deletes object.

```go
func ExampleOperator_Object_delete() {
	acc, _ := newS3Accessor()
	op := NewOperatorFromAccessor(acc)
	object := op.Object("test")
	_ = object.Delete(context.TODO())
}
```

#### List current directory

It returns a [interfaces.ObjectStream](./interfaces/stream.go).

```go
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
```

### Layers
#### Retry
```go
func ExampleOperator_Layer_retry() {
	acc, _ := newS3Accessor()
	op := NewOperatorFromAccessor(acc)

	seed := time.Now().UnixNano()
	random := rand.New(rand.NewSource(seed))

	// retry layer
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
```

See more backoff [strategies](https://github.com/Rican7/retry)


#### Logging

```go
func ExampleOperator_Layer_logging() {
	acc, _ := newS3Accessor()
	op := NewOperatorFromAccessor(acc)

	// logger
	logger, _ := zap.NewProduction()
	s := logger.Sugar()
	
	// logging layer
	loggingLayer := layers.NewLoggingLayer(
		layers.SetLogger(
			layers.NewLo
			ggerAdapter(s.Info, s.Infof),
		),
	)

	op.Layer(loggingLayer)
}
```



## License

The Project is licensed under the [Apache License, Version 2.0](./LICENSE).
