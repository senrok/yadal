package yadal

import (
	"context"
	"fmt"
	"github.com/senrok/yadal/providers/s3"
	"os"
	"testing"
)

var (
	DAL_BUCKET            string
	DAL_ENDPOINT          string
	DAL_ACCESS_KEY_ID     string
	DAL_SECRET_ACCESS_KEY string
)

func TestMain(m *testing.M) {
	provider := os.Getenv("DAL_PROVIDER")
	if provider == "" {
		panic(fmt.Errorf("please set the DAL_PROVIDER env"))
	}
	os.Exit(m.Run())
}

func ExampleNewOperatorFromAccessor() {
	acc, _ := s3.NewDriver(context.Background(), s3.Options{
		Bucket:    "",
		Endpoint:  "",
		Root:      "",
		Region:    "",
		AccessKey: "",
		SecretKey: "",
	})
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
