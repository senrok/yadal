package main

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
