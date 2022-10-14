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
<img src="https://github.com/senrok/yadal/actions/workflows/service_test_s3.yml/badge.svg"/>
</p>


**Y**et **A**nother **D**ata **A**ccess **L**ayer: Access data freely, efficiently, and without the tears 😢

inspired by [Databend's OpenDAL](https://github.com/datafuselabs/opendal)

## Table of contents

- [Features](#features)
- [Installation](#installation)
- [Get started](#get-started)
- [Documentation](#documentation)
- [License](#license)

## Features

Access data freely
- [x] Access different storage services in the same way
- [ ] Behavior tests for all services
  - [x] S3 and S3 compatible services
  - [ ] fs: POSIX compatible filesystem

Access data without  the tears 😢
- [x] Powerful Layer middleware
- [ ] Automatic retry support
- [ ] Logging Layer
- [ ] Tracing Layer
- [ ] Metrics Layer
- [ ] Native decompress support
- [ ] Native service-side encryption support

Access data efficiently
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

## License

The Project is licensed under the [Apache License, Version 2.0](./LICENSE).
