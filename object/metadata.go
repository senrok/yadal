package object

import (
	"github.com/senrok/yadal/constants"
	"github.com/senrok/yadal/interfaces"
	"github.com/senrok/yadal/utils"
	"io/fs"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Metadata struct {
	interfaces.ObjectMode
	contentLength *uint64
	contentMD5    *string
	lastModified  *time.Time
	etag          *string
}

func (m Metadata) Mode() interfaces.ObjectMode {
	return m.ObjectMode
}

func (m Metadata) ContentLength() *uint64 {
	return m.contentLength
}

func (m Metadata) ContentMD5() *string {
	return m.contentMD5
}

func (m Metadata) LastModified() *time.Time {
	return m.lastModified
}

func (m Metadata) ETag() *string {
	return m.etag
}

type MetadataOptions = func(metadata *Metadata) error

func NewMetadata(opts ...MetadataOptions) (interfaces.ObjectMetadata, error) {
	metadata := &Metadata{}
	for _, opt := range opts {
		if err := opt(metadata); err != nil {
			return nil, err
		}
	}
	return metadata, nil
}

func SetFromFileInfo(info fs.FileInfo) MetadataOptions {
	return func(metadata *Metadata) error {
		mode := interfaces.Unknown
		if info.Mode().IsRegular() {
			mode = interfaces.FILE
		} else if info.Mode().IsDir() {
			mode = interfaces.DIR
		}
		metadata.ObjectMode = mode
		return SetMetadata(
			uint64(info.Size()),
			info.ModTime(),
			"")(metadata)
	}
}

func SetMode(mode interfaces.ObjectMode) MetadataOptions {
	return func(metadata *Metadata) error {
		metadata.ObjectMode = mode
		return nil
	}
}

func SetMetadata(size uint64, lm time.Time, etag string) MetadataOptions {
	return func(metadata *Metadata) error {
		metadata.etag = &etag
		metadata.lastModified = &lm
		metadata.contentLength = &size
		md5 := strings.Trim(etag, "\"")
		metadata.contentMD5 = &md5
		return nil
	}
}

func ParseContentLength(header http.Header) MetadataOptions {
	return func(metadata *Metadata) error {
		length := header.Get(constants.ContentLength)
		if length != "" {
			l, err := strconv.ParseUint(length, 10, 64)
			if err != nil {
				return err
			}
			metadata.contentLength = &l
		}
		return nil
	}
}

func ParseETag(header http.Header) MetadataOptions {
	return func(metadata *Metadata) error {
		etagStr := header.Get(constants.ETag)
		if etagStr != "" {
			md5 := strings.Trim(etagStr, "\"")
			metadata.etag = &etagStr
			metadata.contentMD5 = &md5
		}
		return nil
	}
}

func ParseLastModified(header http.Header) MetadataOptions {
	return func(metadata *Metadata) error {
		lastModifiedStr := header.Get(constants.LastModified)
		if lastModifiedStr != "" {
			t, err := utils.ParseRFC7231Time(lastModifiedStr)
			if err != nil {
				return err
			}
			metadata.lastModified = &t
		}
		return nil
	}
}

var metadataOptions = []func(header http.Header) MetadataOptions{
	ParseContentLength,
	ParseETag,
	ParseLastModified,
}

func SetMetadataFromHeader(header http.Header) MetadataOptions {
	return func(metadata *Metadata) error {
		for _, option := range metadataOptions {
			if err := option(header)(metadata); err != nil {
				return err
			}
		}
		return nil
	}
}
