package fs

import (
	"context"
	"github.com/senrok/yadal/interfaces"
	"github.com/senrok/yadal/object"
	"os"
)

type DirStream struct {
	*Driver
	root    string
	path    string
	entries []os.DirEntry
}

func (d *DirStream) HasNext() bool {
	return len(d.entries) > 0
}

func (d *DirStream) Next(ctx context.Context) (entry interfaces.Entry, err error) {
	var e os.DirEntry
	e, d.entries = d.entries[0], d.entries[1:]
	info, err := e.Info()
	if err != nil {
		return nil, err
	}
	meta, err := object.NewMetadata(object.SetFromFileInfo(info))
	if err != nil {
		return nil, err
	}
	return object.NewEntry(d.Driver, d.path, meta, false), nil
}
