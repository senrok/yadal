package s3

import (
	"context"
	"encoding/xml"
	"github.com/senrok/yadal/errors"
	"github.com/senrok/yadal/interfaces"
	"github.com/senrok/yadal/object"
	"github.com/senrok/yadal/utils"
	"net/http"
	"strings"
	"time"
)

type DirStream struct {
	*Driver
	root  string
	path  string
	token string

	done bool
}

func (d *DirStream) NextPage(ctx context.Context) ([]interfaces.Entry, error) {
	if d.done {
		return nil, nil
	}
	resp, err := d.ListObjects(ctx, d.path, d.token)
	if err != nil {
		return nil, errors.Wrap(errors.ErrListFailed, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.ParseS3Error(errors.ErrListFailed, d.path, resp)
	}
	output := Output{}
	err = xml.NewDecoder(resp.Body).Decode(&output)
	if err != nil {
		return nil, errors.Wrap(errors.ErrListFailed, err)
	}

	if output.IsTruncated != nil {
		d.done = !*output.IsTruncated
	} else if output.NextContinuationToken != nil {
		d.done = *output.NextContinuationToken == ""
	} else {
		d.done = len(output.CommonPrefixes) == 0 && len(output.Contents) == 0
	}
	if output.NextContinuationToken != nil {
		d.token = *output.NextContinuationToken
	}

	entries := make([]interfaces.Entry, 0, len(output.CommonPrefixes)+len(output.Contents))
	for _, prefix := range output.CommonPrefixes {
		path, err := utils.BuildRealPath(d.root, prefix.Prefix)
		if err != nil {
			return nil, err
		}

		entries = append(entries,
			object.NewEntry(
				d.Driver,
				path,
				object.Metadata{
					ObjectMode: interfaces.DIR,
				},
				false),
		)
	}

	for _, content := range output.Contents {
		// s3 could return the dir itself in contents
		// which endswith `/`.
		// We should ignore them.
		if strings.HasSuffix(content.Key, "/") {
			continue
		}
		meta, err := object.NewMetadata(object.SetMode(interfaces.FILE), object.SetMetadata(content.Size, content.LastModified, content.ETag))
		if err != nil {
			return nil, err
		}
		path, err := utils.BuildRealPath(d.root, content.Key)
		if err != nil {
			return nil, err
		}
		entries = append(entries,
			object.NewEntry(
				d.Driver,
				path,
				meta,
				false),
		)
	}
	return entries, nil
}

func NewDirStream(d *Driver, root, path string) interfaces.ObjectPageStream {
	return &DirStream{
		Driver: d,
		root:   root,
		path:   path,
		token:  "",
		done:   false,
	}
}

type Output struct {
	IsTruncated           *bool                `xml:"IsTruncated"`
	NextContinuationToken *string              `xml:"NextContinuationToken"`
	CommonPrefixes        []OutputCommonPrefix `xml:"CommonPrefixes"`
	Contents              []OutputContent      `xml:"Contents"`
}

type OutputContent struct {
	Key          string    `xml:"Key"`
	Size         uint64    `xml:"size"`
	LastModified time.Time `xml:"LastModified"`
	ETag         string    `xml:"ETag"`
}

type OutputCommonPrefix struct {
	Prefix string `xml:"Prefix"`
}
