package s3

import (
	"context"
	"fmt"
	"github.com/senrok/yadal/interfaces"
	"github.com/senrok/yadal/options"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"os"
	"strings"
	"testing"
)

var (
	DAL_S3_BUCKET            = os.Getenv("DAL_S3_BUCKET")
	DAL_S3_ENDPOINT          = os.Getenv("DAL_S3_ENDPOINT")
	DAL_S3_ACCESS_KEY_ID     = os.Getenv("DAL_S3_ACCESS_KEY_ID")
	DAL_S3_SECRET_ACCESS_KEY = os.Getenv("DAL_S3_SECRET_ACCESS_KEY")
	DAL_S3_TEST              = os.Getenv("DAL_S3_TEST")
)

func TestMain(m *testing.M) {
	if DAL_S3_TEST == "on" {
		if DAL_S3_BUCKET == "" || DAL_S3_ENDPOINT == "" {
			panic(fmt.Errorf("please set the DAL_S3_BUCKET and DAL_S3_ENDPOINT env"))
		}
		if DAL_S3_ACCESS_KEY_ID == "" {
			log.Printf("Reading the empty DAL_S3_ACCESS_KEY_ID\n")
		}
		if DAL_S3_SECRET_ACCESS_KEY == "" {
			log.Printf("Reading the empty DAL_S3_SECRET_ACCESS_KEY\n")
		}
	}

	os.Exit(m.Run())
}

func setupDriver(t *testing.T) interfaces.Accessor {
	d, err := NewDriver(context.Background(), Options{
		Bucket:                     DAL_S3_BUCKET,
		Endpoint:                   DAL_S3_ENDPOINT,
		Root:                       "",
		Region:                     "",
		AccessKey:                  DAL_S3_ACCESS_KEY_ID,
		SecretKey:                  DAL_S3_SECRET_ACCESS_KEY,
		SSEncryption:               nil,
		SSEncryptionAwsKmsKeyId:    nil,
		SSEncryptionCustomerAlgo:   nil,
		SSEncryptionCustomerKey:    nil,
		SSEncryptionCustomerKeyMD5: nil,
	})
	assert.Nil(t, err)
	assert.NotNil(t, d)
	return d
}

func TestNewDriver(t *testing.T) {
	d, err := NewDriver(context.Background(), Options{
		Bucket:                     DAL_S3_BUCKET,
		Endpoint:                   DAL_S3_ENDPOINT,
		Root:                       "/",
		Region:                     "",
		AccessKey:                  DAL_S3_ACCESS_KEY_ID,
		SecretKey:                  DAL_S3_SECRET_ACCESS_KEY,
		SSEncryption:               nil,
		SSEncryptionAwsKmsKeyId:    nil,
		SSEncryptionCustomerAlgo:   nil,
		SSEncryptionCustomerKey:    nil,
		SSEncryptionCustomerKeyMD5: nil,
	})
	assert.Nilf(t, err, "%s", err)
	assert.NotNil(t, d)
}

func TestDriver_Create(t *testing.T) {
	acc := setupDriver(t)
	err := acc.Create(context.Background(), "test-dir/", options.CreateOptions{})
	assert.Nilf(t, err, "%s", err)
}

func TestDriver_Write(t *testing.T) {
	acc := setupDriver(t)
	text := "hello world"
	body := strings.NewReader(text)

	size, err := acc.Write(
		context.Background(),
		"test-dir/hello.txt",
		options.WriteOptions{Size: uint64(len(text))},
		body,
	)
	assert.Equal(t, uint64(len(text)), size)
	assert.Nilf(t, err, "%s", err)
}

func TestDriver_Read(t *testing.T) {
	acc := setupDriver(t)
	read, err := acc.Read(context.Background(), "test-dir/hello.txt", options.ReadOptions{})
	assert.Nilf(t, err, "%s", err)
	b, _ := io.ReadAll(read)
	assert.Equal(t, "hello world", string(b))
}

func TestDriver_List(t *testing.T) {
	acc := setupDriver(t)
	s, err := acc.List(context.Background(), "test-dir/", options.ListOptions{})
	assert.Nilf(t, err, "%s", err)
	assert.True(t, s.HasNext())
	e, err := s.Next(context.Background())
	assert.Nilf(t, err, "%s", err)
	assert.Equal(t, "test-dir/hello.txt", e.Path())
	assert.False(t, s.HasNext())
}

func TestDriver_Delete(t *testing.T) {
	acc := setupDriver(t)
	err := acc.Delete(context.Background(), "test-dir/hello.txt", options.DeleteOptions{})
	assert.Nilf(t, err, "%s", err)
}
