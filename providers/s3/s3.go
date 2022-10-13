package s3

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/senrok/yadal/constants"
	"github.com/senrok/yadal/errors"
	"github.com/senrok/yadal/interfaces"
	"github.com/senrok/yadal/logger"
	"github.com/senrok/yadal/object"
	"github.com/senrok/yadal/options"
	"github.com/senrok/yadal/providers"
	"github.com/senrok/yadal/utils"
	"io"
	"net/http"
	"strings"
)

type Driver struct {
	bucket                     string
	endpoint                   string
	region                     string
	root                       string
	client                     http.Client
	signer                     Signer
	SSEncryption               *string
	SSEncryptionAwsKmsKeyId    *string
	SSEncryptionCustomerAlgo   *string
	SSEncryptionCustomerKey    *string
	SSEncryptionCustomerKeyMD5 *string
	logger.Logger
}

func (d *Driver) Metadata() interfaces.Metadata {
	return providers.NewMetadata(interfaces.S3, d.root, d.bucket, interfaces.Read|interfaces.Write|interfaces.List|interfaces.PreSign|interfaces.Multipart)
}

func (d *Driver) Create(ctx context.Context, path string, _ options.CreateOptions) error {
	req, err := d.putObjectRequest(ctx, path, nil, nil)
	if err != nil {
		return errors.Wrap(errors.ErrCreateFailed, err)
	}

	if err = d.signer.Sign(req, nil); err != nil {
		return errors.Wrap(errors.ErrCreateFailed, err)
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return errors.Wrap(errors.ErrCreateFailed, err)
	}

	switch resp.StatusCode {
	case http.StatusCreated, http.StatusOK:
		return nil
	default:
		return errors.ParserError(errors.ErrCreateFailed, path, resp)
	}
}

func (d *Driver) Read(ctx context.Context, path string, args options.ReadOptions) (io.ReadCloser, error) {
	resp, err := d.GetObject(ctx, path, args.Offset, args.Size)
	if err != nil {
		return nil, errors.Wrap(errors.ErrReadFailed, err)
	}
	switch resp.StatusCode {
	case http.StatusCreated, http.StatusOK:
		return resp.Body, nil
	default:
		return nil, errors.ParserError(errors.ErrReadFailed, path, resp)
	}
}

func (d *Driver) Write(ctx context.Context, path string, args options.WriteOptions, reader io.ReadSeeker) (uint64, error) {
	req, err := d.putObjectRequest(ctx, path, &args.Size, reader)
	if err != nil {
		return 0, errors.Wrap(errors.ErrWriteFailed, err)
	}
	if err = d.signer.Sign(req, reader); err != nil {
		return 0, errors.Wrap(errors.ErrWriteFailed, err)
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return 0, errors.Wrap(errors.ErrWriteFailed, err)
	}

	switch resp.StatusCode {
	case http.StatusCreated, http.StatusOK:
		return args.Size, nil
	default:
		return 0, errors.ParserError(errors.ErrWriteFailed, path, resp)
	}
}

func (d *Driver) Stat(ctx context.Context, path string, args options.StatOptions) (interfaces.ObjectMetadata, error) {
	if path == "/" {
		return object.Metadata{ObjectMode: interfaces.DIR}, nil
	}

	resp, err := d.HeadObject(ctx, path)
	if err != nil {
		return nil, errors.Wrap(errors.ErrStatFailed, err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		md, err := object.NewMetadata(
			object.SetMode(interfaces.ObjectModeFromPath(path)),
			object.SetMetadataFromHeader(resp.Header),
		)
		if err != nil {
			return nil, errors.Wrap(errors.ErrStatFailed, err)
		}
		return md, nil

	case http.StatusNotFound:
		if strings.HasPrefix(path, "/") {
			return object.Metadata{ObjectMode: interfaces.DIR}, nil
		}
		fallthrough // handles other cases
	default:
		return nil, errors.ParserError(errors.ErrStatFailed, path, resp)
	}
}

func (d *Driver) Delete(ctx context.Context, path string, args options.DeleteOptions) error {
	resp, err := d.DeleteObject(ctx, path)
	if err != nil {
		return errors.Wrap(errors.ErrDeleteFailed, err)
	}
	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	default:
		return errors.ParserError(errors.ErrDeleteFailed, path, resp)
	}
}

func (d *Driver) List(ctx context.Context, path string, args options.ListOptions) (interfaces.ObjectStream, error) {
	return object.NewObjectStream(&DirStream{
		Driver: d,
		root:   d.root,
		path:   path,
	}), nil
}

func (d *Driver) PreSign(ctx context.Context, path string, args options.PreSignOptions) (req *http.Request, err error) {

	switch args.Op {
	case options.ReadOp:
		if req, err = d.getObjectRequest(ctx, path, args.Offset, args.ReadOptions.Size); err != nil {
			return nil, errors.Wrap(errors.ErrPreSignFailed, err)
		}
	case options.WriteOp:
		if req, err = d.putObjectRequest(ctx, path, nil, nil); err != nil {
			return nil, errors.Wrap(errors.ErrPreSignFailed, err)
		}
	case options.WriteMultipartOp:
		if req, err = d.S3UploadPartRequest(ctx, path, args.UploadId, args.PartNumber, nil, nil); err != nil {
			return nil, errors.Wrap(errors.ErrPreSignFailed, err)
		}
	default:
		return nil, errors.Wrap(errors.ErrPreSignFailed, errors.ErrUnknownPreSignOperation)
	}

	if err = d.signer.Sign(req, nil); err != nil {
		return nil, errors.Wrap(errors.ErrPreSignFailed, err)
	}
	return
}

func (d *Driver) CreateMultipart(ctx context.Context, path string, args options.CreateMultipart) (string, error) {
	resp, err := d.S3InitiateMultipartUpload(ctx, path)
	if err != nil {
		return "", errors.Wrap(errors.ErrCreateMultipartFailed, err)
	}
	switch resp.StatusCode {
	case http.StatusOK:
		output := InitiateMultipartUploadResult{}
		if err = xml.NewDecoder(resp.Body).Decode(&output); err != nil {
			return "", errors.Wrap(errors.ErrCreateMultipartFailed, err)
		}
		return output.UploadId, nil
	default:
		return "", errors.ParserError(errors.ErrCreateMultipartFailed, path, resp)
	}
}

func (d *Driver) WriteMultipart(ctx context.Context, path string, args options.WriteMultipart, reader io.ReadSeeker) (interfaces.ObjectPart, error) {
	req, err := d.S3UploadPartRequest(ctx, path, args.UploadId, args.PartNumber, &args.Size, io.NopCloser(reader))
	if err != nil {
		return nil, errors.Wrap(errors.ErrWriteMultipartFailed, err)
	}
	if err = d.signer.Sign(req, reader); err != nil {
		return nil, err
	}
	resp, err := d.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(errors.ErrWriteMultipartFailed, err)
	}
	switch resp.StatusCode {
	case http.StatusOK:
		return object.ObjectPart{
			PartNumber: args.PartNumber,
			ETag:       resp.Header.Get(constants.ETag),
		}, nil

	default:
		return nil, errors.ParserError(errors.ErrWriteMultipartFailed, path, resp)
	}
}

func (d *Driver) CompleteMultipart(ctx context.Context, path string, args options.CompleteMultipart) error {
	// TODO: any other way?
	parts := make([]interfaces.ObjectPart, 0, len(args.ObjectParts))
	for _, part := range args.ObjectParts {
		parts = append(parts, part)
	}
	resp, err := d.S3CompleteMultipartUpload(ctx, path, args.UploadId, parts)
	if err != nil {
		return errors.Wrap(errors.ErrCompleteMultipartFailed, err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	default:
		return errors.ParserError(errors.ErrCompleteMultipartFailed, path, resp)
	}
}

func (d *Driver) AbortMultipart(ctx context.Context, path string, args options.AbortMultipart) error {
	resp, err := d.S3AbortMultipartUpload(ctx, path, args.UploadId)
	if err != nil {
		return errors.Wrap(errors.ErrAbortMultipartFailed, err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	default:
		return errors.ParserError(errors.ErrAbortMultipartFailed, path, resp)
	}
}

func (d *Driver) detectRegion(ctx context.Context, bucket string) (endpoint string, region string, err error) {
	endpoint = d.endpoint
	if !strings.HasPrefix(endpoint, "http") {
		endpoint = fmt.Sprintf("https://%s", endpoint)
	}

	endpoint = strings.Replace(endpoint, fmt.Sprintf("//%s.", d.bucket), "//", 1)

	if d.region != "" {
		return endpoint, d.region, nil
	}

	url := fmt.Sprintf("%s/%s", endpoint, bucket)
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return "", "", err
	}
	resp, err := d.client.Do(req)
	if err != nil {
		return "", "", err
	}
	switch resp.StatusCode {
	case http.StatusOK, http.StatusForbidden:
		region = resp.Header.Get(constants.XAmzBucketRegion)
		if region == "" {
			region = "us-east-1"
		}
		return
	default:
		return "", "", errors.ParserError(errors.ErrDetectRegionFailed, req.URL.String(), resp)
	}
}

func NewDriver(ctx context.Context, opt Options) (interfaces.Accessor, error) {
	region := opt.Region

	d := &Driver{
		bucket:                     opt.Bucket,
		endpoint:                   opt.Endpoint,
		root:                       utils.NormalizeRoot(opt.Root),
		client:                     http.Client{},
		region:                     region,
		SSEncryption:               opt.SSEncryption,
		SSEncryptionAwsKmsKeyId:    opt.SSEncryptionAwsKmsKeyId,
		SSEncryptionCustomerAlgo:   opt.SSEncryptionCustomerAlgo,
		SSEncryptionCustomerKey:    opt.SSEncryptionCustomerKey,
		SSEncryptionCustomerKeyMD5: opt.SSEncryptionCustomerKeyMD5,
	}

	if opt.Region == "" {
		var err error
		if d.endpoint, d.region, err = d.detectRegion(ctx, d.bucket); err != nil {
			return nil, err
		}
	}

	if opt.EnableVirtualHostStyle {
		d.endpoint = strings.Replace(d.endpoint, "//", fmt.Sprintf("//%s", opt.Bucket), 1)
	} else {
		d.endpoint = d.endpoint + "/" + opt.Bucket
	}

	d.signer = NewSigner("s3", d.region, opt.AccessKey, opt.SecretKey, true)

	return d, nil
}

type InitiateMultipartUploadResult struct {
	UploadId string `xml:"UploadId"`
}
