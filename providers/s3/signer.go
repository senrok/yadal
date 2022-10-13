// Package aws4 signs HTTP requests as prescribed in
// http://docs.amazonwebservices.com/general/latest/gr/signature-version-4.html
// from: github.com/bmizerany/aws4

package s3

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"io"
	"net/http"
	"time"
)

type Signer interface {
	Sign(r *http.Request, reader io.ReadSeeker) error
}

type signer struct {
	*v4.Signer
	Anonymous bool
	Service   string
	Region    string
}

func (s signer) Sign(r *http.Request, reader io.ReadSeeker) error {
	if s.Anonymous {
		return nil
	}
	_, err := s.Signer.Sign(r, reader, s.Service, s.Region, time.Now())
	if err != nil {
		return err
	}
	return nil
}

type naiveSignerProvider struct {
	AccessKey string
	SecretKey string
}

func (s naiveSignerProvider) Retrieve() (credentials.Value, error) {
	return credentials.Value{
		AccessKeyID:     s.AccessKey,
		SecretAccessKey: s.SecretKey,
	}, nil
}

func (s naiveSignerProvider) IsExpired() bool {
	return false
}

func NewNaiveSignerProvider(a, s string) credentials.Provider {
	return &naiveSignerProvider{
		AccessKey: a,
		SecretKey: s,
	}
}

func NewSigner(service, region, accessKey, secretKey string, allow_anonymous bool) Signer {
	return &signer{
		Anonymous: allow_anonymous && accessKey == "" && secretKey == "",
		Signer:    v4.NewSigner(credentials.NewCredentials(NewNaiveSignerProvider(accessKey, secretKey))),
		Service:   service,
		Region:    region,
	}
}
