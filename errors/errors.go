package errors

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrUnsupportedMethod = errors.New("unsupported method")

	ErrCreateFailed  = errors.New("create operation failed")
	ErrReadFailed    = errors.New("read operation failed")
	ErrWriteFailed   = errors.New("write operation failed")
	ErrStatFailed    = errors.New("stat operation failed")
	ErrDeleteFailed  = errors.New("delete operation failed")
	ErrPreSignFailed = errors.New("presign operation failed")

	ErrCreateMultipartFailed   = errors.New("create multipart operation failed")
	ErrWriteMultipartFailed    = errors.New("write multipart operation failed")
	ErrCompleteMultipartFailed = errors.New("complete multipart operation failed")
	ErrAbortMultipartFailed    = errors.New("abort multipart operation failed")

	ErrUnknownPreSignOperation = errors.New("unknown presign operation")

	ErrDetectRegionFailed = errors.New("detect region failed")

	ErrListFailed = errors.New("list operation failed")

	ErrNotFound         = errors.New("not found")
	ErrPermissionDenied = errors.New("permission denied")
	ErrInterrupted      = errors.New("err interrupted")
	ErrOther            = errors.New("unknown error")
)

type ObjectError struct {
	source error
	kind   error
	path   string
	body   []byte
}

func (error ObjectError) Is(other error) bool {
	return errors.Is(error.source, other) || errors.Is(error.kind, other)
}

func (error ObjectError) Error() string {
	return fmt.Sprintf("kind %s\nsource:%s\npath: %s\nbody: %s\n", error.kind, error.source, error.path, string(error.body))
}

func ParserError(err error, path string, resp *http.Response) error {
	var kind error
	switch resp.StatusCode {
	case http.StatusNotFound:
		kind = ErrNotFound
	case http.StatusForbidden:
		kind = ErrPermissionDenied
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		kind = ErrInterrupted
	default:
		kind = ErrOther
	}
	b, _ := io.ReadAll(resp.Body)
	return &ObjectError{
		source: err,
		kind:   kind,
		path:   path,
		body:   b,
	}
}

func Wrap(err error, child error) error {
	return fmt.Errorf("%w\ndue:%s", err, child)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}
