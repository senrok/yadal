package s3

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"github.com/senrok/yadal/constants"
	"github.com/senrok/yadal/interfaces"
	"github.com/senrok/yadal/options"
	"github.com/senrok/yadal/utils"
	"io"
	"net/http"
	"strconv"
)

func (d *Driver) insertSseHeaders(req *http.Request, isWrite bool) {
	if isWrite {
		if d.SSEncryption != nil {
			req.Header.Set(constants.XAmzServerSideEncryption, *d.SSEncryption)
		}
		if d.SSEncryptionAwsKmsKeyId != nil {
			req.Header.Set(constants.XAmzServerSideEncryptionAwsKmsKeyId, *d.SSEncryptionAwsKmsKeyId)
		}
	}
	if d.SSEncryptionCustomerAlgo != nil {
		req.Header.Set(constants.XAmzServerSideEncryptionCustomerAlgorithm, *d.SSEncryptionCustomerAlgo)
	}
	if d.SSEncryptionCustomerKey != nil {
		req.Header.Set(constants.XAmzServerSideEncryptionCustomerKey, *d.SSEncryptionCustomerKey)
	}
	if d.SSEncryptionCustomerKeyMD5 != nil {
		req.Header.Set(constants.XAmzServerSideEncryptionCustomerKeyMd5, *d.SSEncryptionCustomerKeyMD5)
	}
}

func (d *Driver) buildUrl(path string) (string, error) {
	p, err := utils.BuildAbsPath(d.root, path)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/%s", d.endpoint, utils.EncodePath(p))

	return url, nil
}

func (d *Driver) getObjectRequest(_ context.Context, path string, offset, size *uint64) (*http.Request, error) {
	url, err := d.buildUrl(path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if (offset != nil && *offset != 0) || size != nil {
		req.Header.Set("range", options.NewBytesRange(offset, size).String())
	}

	// SSE headers
	d.insertSseHeaders(req, false)

	return req, nil
}

func (d *Driver) GetObject(ctx context.Context, path string, offset, size *uint64) (*http.Response, error) {
	req, err := d.getObjectRequest(ctx, path, offset, size)
	if err != nil {
		return nil, err
	}

	if err = d.signer.Sign(req, nil); err != nil {
		return nil, err
	}

	return d.client.Do(req)
}

func (d *Driver) putObjectRequest(_ context.Context, path string, size *uint64, body io.ReadSeeker) (*http.Request, error) {
	url, err := d.buildUrl(path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}

	if size != nil {
		req.Header.Set(constants.ContentLength, fmt.Sprintf("%d", int64(*size)))
		req.ContentLength = int64(*size)
	}

	// SSE headers
	d.insertSseHeaders(req, true)

	return req, nil
}

func (d *Driver) HeadObject(_ context.Context, path string) (*http.Response, error) {
	url, err := d.buildUrl(path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return nil, err
	}

	// SSE headers
	d.insertSseHeaders(req, false)

	err = d.signer.Sign(req, nil)
	if err != nil {
		return nil, err
	}

	return d.client.Do(req)
}

func (d *Driver) DeleteObject(_ context.Context, path string) (*http.Response, error) {
	url, err := d.buildUrl(path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	if err = d.signer.Sign(req, nil); err != nil {
		return nil, err
	}

	return d.client.Do(req)
}

func (d *Driver) ListObjects(_ context.Context, path string, continuationToken string) (*http.Response, error) {
	p, err := utils.BuildAbsPath(d.root, path)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s?list-type=2&delimiter=/&prefix=%s", d.endpoint, utils.EncodePath(p))

	if continuationToken != "" {
		url += fmt.Sprintf("&continuation-token=%s", continuationToken)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if err = d.signer.Sign(req, nil); err != nil {
		return nil, err
	}

	return d.client.Do(req)
}

func (d *Driver) S3InitiateMultipartUpload(_ context.Context, path string) (*http.Response, error) {
	p, err := utils.BuildAbsPath(d.root, path)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/%s?uploads", d.endpoint, utils.EncodePath(p))

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}

	if err = d.signer.Sign(req, nil); err != nil {
		return nil, err
	}

	return d.client.Do(req)
}

func (d *Driver) S3UploadPartRequest(ctx context.Context, path string, uploadId string, partNumber uint, size *uint64, body io.ReadCloser) (*http.Request, error) {
	p, err := utils.BuildAbsPath(d.root, path)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/%s?partNumber=%d&uploadId=%s", d.endpoint, utils.EncodePath(p), partNumber, uploadId)

	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}
	if size != nil {
		req.Header.Set(constants.ContentLength, strconv.FormatUint(*size, 10))
	}

	return req, nil
}

func (d *Driver) S3CompleteMultipartUpload(_ context.Context, path, uploadId string, parts []interfaces.ObjectPart) (*http.Response, error) {
	p, err := utils.BuildAbsPath(d.root, path)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/%s?uploadId=%s", d.endpoint, utils.EncodePath(p), uploadId)

	body, err := xml.Marshal(NewCompleteMultipartUploadFromObjectParts(parts))
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, io.NopCloser(bytes.NewReader(body)))
	if err != nil {
		return nil, err
	}

	req.Header.Set(constants.ContentLength, strconv.FormatUint(uint64(len(body)), 10))

	req.Header.Set(constants.ContentType, "application/xml")

	if err = d.signer.Sign(req, nil); err != nil {
		return nil, err
	}

	return d.client.Do(req)
}

func (d *Driver) S3AbortMultipartUpload(_ context.Context, path, uploadId string) (*http.Response, error) {
	p, err := utils.BuildAbsPath(d.root, path)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/%s?uploadId=%s", d.endpoint, utils.EncodePath(p), uploadId)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	if err = d.signer.Sign(req, nil); err != nil {
		return nil, err
	}

	return d.client.Do(req)
}
